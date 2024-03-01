package idcache_test

import (
	"testing"

	"github.com/hoodnoah/eve_market/monitor/idcache"
)

func Test_IDCache_Label(t *testing.T) {
	t.Run("it should return the expected result for a result already cached", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		cache.SetKnownIDs(map[int]string{1: "hello world"})

		res, err := testSetup.Cache.Label(1)
		if err != nil {
			t.Fatalf("expected to receive a successful result: %s", err)
		}

		if res != "hello world" {
			t.Fatalf("expected to receive string 'hello world', received '%s'", res)
		}
	})

	t.Run("it should reach out to the API for a result not yet cached", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		// type assertion to narrow down to FakeClient from IRequestClient
		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		fakeClient.SetResponseMap(map[int]string{1: "hello world"})
		_, _ = cache.Label(1)
		args := fakeClient.Args

		if len(args) != 1 {
			t.Fatalf("expected only 1 request from the cache, received %d", len(args))
		}

		arg := args[0][0]

		if arg != 1 {
			t.Fatalf("expected to call client with 1, received %d", arg)
		}
	})

	t.Run("it should cache the result when it has to fetch a result not yet cached", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		testClient := testSetup.Client
		var fakeClient *FakeClient

		if fc, ok := testClient.(*FakeClient); ok {
			fakeClient = fc
		}

		fakeClient.SetResponseMap(map[int]string{1: "hello world"})

		_, _ = cache.Label(1)
		_, _ = cache.Label(1)

		args := fakeClient.Args

		if len(args) != 1 {
			t.Fatalf("expected to have only hit the client once, received %d hits", len(args))
		}
	})

	t.Run("it should return the provided result faithfully when it has to hit the API", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		fakeClient.SetResponseMap(map[int]string{1: "hello world"})

		val, err := cache.Label(1)

		if err != nil {
			t.Fatalf("expected a successful result, received: %s", err)
		}

		if val != "hello world" {
			t.Fatalf("expected value to be 'hello world', received '%s'", val)
		}
	})

	t.Run("it should come up with a sensible name when an ID is deemed invalid by the API", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		fakeClient.SetInvalidIDS([]int{1})

		val, err := cache.Label(1)

		if err != nil {
			t.Fatalf("should not fail on an invalid id, but did so: %s", err)
		}

		if val != "invalidID_1" {
			t.Fatalf("expected invalidID_1 as response for invalid id, received %s", val)
		}
	})

	t.Run("it should not come up with an ID when the request fails for any other reason", func(t *testing.T) {
		testSetup := setup()
		cache := testSetup.Cache
		var fakeClient *FakeClient

		if fc, ok := testSetup.Client.(*FakeClient); ok {
			fakeClient = fc
		}

		fakeClient.SetFail(true, idcache.InvalidRequestError)

		val, err := cache.Label(1)
		if val != "" {
			t.Fatalf("expected cache to return empty string, received %s", val)
		}

		if err == nil {
			t.Fatal("expected cache to fail, but the error was nil")
		}
	})
}
