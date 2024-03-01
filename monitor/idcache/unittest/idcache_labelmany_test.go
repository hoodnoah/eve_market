package idcache_test

import (
	"fmt"
	"testing"
)

func Test_IDCache_LabelMany(t *testing.T) {

	t.Run("it should label several ids when they all exist in the cache from the start", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache

		labeledIds := map[int]string{
			1: "hello",
			2: "world",
			3: "this",
			4: "is",
			5: "a",
			6: "unit",
			7: "test",
		}

		idsToLabel := []int{1, 2, 3, 4, 5, 6, 7}

		cache.SetKnownIDs(labeledIds)

		result, err := cache.LabelMany(idsToLabel)
		if err != nil {
			t.Fatalf("expected cache to label items that have been cached, instead received error: %s", err)
		}

		for _, id := range idsToLabel {
			if result[id] != labeledIds[id] {
				t.Fatalf("expected cache to return same values as provided; id %d was supposed to resolve to %s, but actually resolved to %s", id, labeledIds[id], result[id])
			}
		}
	})

	t.Run("it should not hit the API when all the ids are already cached", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		labeledIds := map[int]string{
			1: "hello",
			2: "world",
			3: "this",
			4: "is",
			5: "a",
			6: "unit",
			7: "test",
		}

		idsToLabel := []int{1, 2, 3, 4, 5, 6, 7}

		cache.SetKnownIDs(labeledIds)

		_, err := cache.LabelMany(idsToLabel)
		if err != nil {
			t.Fatalf("expected no error, instead received error: %s", err)
		}

		if len(fakeClient.Args) != 0 {
			t.Fatalf("expected no calls to be made to the http client, instead it was called %d times", len(fakeClient.Args))
		}
	})

	t.Run("it should hit the API when the requested ids are not cached", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		labeledIds := map[int]string{
			1: "hello",
			2: "world",
			3: "this",
			4: "is",
			5: "a",
			6: "unit",
			7: "test",
		}

		inputIds := []int{1, 2, 3, 4, 5, 6, 7}

		fakeClient.SetResponseMap(labeledIds)

		res, err := cache.LabelMany(inputIds)
		if err != nil {
			t.Fatalf("expected to hit the API for results, errored instead: %s", err)
		}

		for id := range res {
			if res[id] != labeledIds[id] {
				t.Fatalf("expected id #%d to return value '%s', received '%s'", id, labeledIds[id], res[id])
			}
		}

		if len(fakeClient.Args) != 1 {
			t.Fatalf("expected only 1 API hit, received %d", len(fakeClient.Args))
		}

		args := fakeClient.Args[0]
		if len(args) != len(inputIds) {
			t.Fatalf("expected the http client to be called with %d ids, received %d", len(inputIds), len(args))
		}

		for index := range args {
			if args[index] != inputIds[index] {
				t.Fatalf("expected the %d element in the argument to the client to be %d, received %d", index, inputIds[index], args[index])
			}
		}
	})

	t.Run("it should only hit the api for non-cached values", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		knownIds := map[int]string{
			1: "hello",
			3: "this",
			4: "is",
			5: "a",
			6: "unit",
			7: "test",
		}

		fakeClient.SetResponseMap(map[int]string{2: "world"})

		idsToLabel := []int{1, 2, 3, 4, 5, 6, 7}

		cache.SetKnownIDs(knownIds)

		_, err := cache.LabelMany(idsToLabel)
		if err != nil {
			t.Fatalf("expected a successful response, received an error: %s", err)
		}

		if len(fakeClient.Args) != 1 {
			t.Fatalf("expected API to be hit once, instead hit it %d times", len(fakeClient.Args))
		}

		args := fakeClient.Args[0]
		if len(args) != 1 && args[0] != 2 {
			t.Fatalf("expected API to be hit for only id 2, hit it for %v", args)
		}
	})

	t.Run("it should return a correct result when some ids are cached, and others aren't", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}
		allIds := map[int]string{
			1: "hello",
			2: "world",
			3: "this",
			4: "is",
			5: "a",
			6: "unit",
			7: "test",
		}

		knownIds := map[int]string{
			1: "hello",
			3: "this",
			6: "unit",
			7: "test",
		}

		fakeClient.SetResponseMap(map[int]string{2: "world", 4: "is", 5: "a"})

		idsToLabel := []int{1, 2, 3, 4, 5, 6, 7}

		cache.SetKnownIDs(knownIds)

		res, err := cache.LabelMany(idsToLabel)
		if err != nil {
			t.Fatalf("expected a successful response, received an error: %s", err)
		}

		for id, val := range res {
			if val != allIds[id] {
				t.Fatalf("expected id #%d to have value '%s', received '%s'", id, allIds[id], val)
			}
		}
	})

	t.Run("it should auto-fill invalid ids, and return the correct ids for those which are valid", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		expectedIds := map[int]string{
			1: "hello",
			2: "world",
			3: "unknownID_3",
			4: "is",
			5: "a",
			6: "unit",
			7: "unknownID_7",
		}

		knownIds := map[int]string{
			1: "hello",
			6: "unit",
		}

		idsToLabel := []int{1, 2, 3, 4, 5, 6, 7}

		cache.SetKnownIDs(knownIds)
		fakeClient.SetInvalidIDS([]int{3, 7})
		fakeClient.SetResponseMap(map[int]string{
			1: "hello",
			2: "world",
			4: "is",
			5: "a",
			6: "unit",
		})

		res, err := cache.LabelMany(idsToLabel)
		if err != nil {
			t.Fatalf("expected success, received error: %s", err)
		}

		for id, val := range expectedIds {
			if res[id] != val {
				t.Fatalf("expected id #%d to have value '%s', received '%s'", id, val, res[id])
			}
		}
	})

	t.Run("it should never request more than a thousand ids in one go", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		idValues := map[int]string{}
		idsToFetch := []int{}
		for i := range 2500 {
			idsToFetch = append(idsToFetch, i)
			idValues[i] = fmt.Sprintf("label_%d", i)
		}

		fakeClient.ResponseMap = idValues

		_, err := cache.LabelMany(idsToFetch)
		if err != nil {
			t.Fatalf("expected a successful response, got err: %s", err)
		}
		args := fakeClient.Args

		for _, arg := range args {
			if len(arg) > 1000 {
				t.Fatalf("expected no more than 1,000 ids per request, received: %d", len(arg))
			}
		}

	})
}
