<script setup lang="ts">
import { computed, ref } from 'vue'

// types
import type { IValueItem } from '@/types/ValueItem'

const searchTerm = ref('')
const possibleItems = ref<IValueItem[]>([])
const emit = defineEmits(['add-selection'])

possibleItems.value = [
  { id: 1, value: 'Apple' },
  { id: 2, value: 'Banana' },
  { id: 3, value: 'Cherry' },
  { id: 4, value: 'Date' },
  { id: 5, value: 'Elderberry' },
  { id: 6, value: 'Fig' },
  { id: 7, value: 'Grape' },
  { id: 8, value: 'Honeydew' },
  { id: 9, value: 'Icaco' },
  { id: 10, value: 'Jackfruit' },
  { id: 11, value: 'Kiwi' },
  { id: 12, value: 'Lemon' },
  { id: 13, value: 'Mango' },
  { id: 14, value: 'Nectarine' },
  { id: 15, value: 'Orange' },
  { id: 16, value: 'Papaya' },
  { id: 17, value: 'Quince' },
  { id: 18, value: 'Raspberry' },
  { id: 19, value: 'Strawberry' },
  { id: 20, value: 'Tangerine' },
  { id: 21, value: 'Ugli' },
  { id: 22, value: 'Vanilla' },
  { id: 23, value: 'Watermelon' },
  { id: 24, value: 'Xigua' },
  { id: 25, value: 'Yuzu' },
  { id: 26, value: 'Zucchini' }
]

const filteredSelections = computed(() => {
  return possibleItems.value.filter((item) => {
    return item.value.toLowerCase().includes(searchTerm.value.toLowerCase())
  })
})

function handleSelectItem(itemId: number) {
  console.log('handleSelectItem: ', itemId)
  const item = possibleItems.value.find((i) => i.id === itemId)
  if (item) {
    console.log('setting searchTerm value to: ', item.value)
    searchTerm.value = item.value
  }
}

function handleAddItem() {
  const item = possibleItems.value.find((i) => i.value === searchTerm.value)
  if (item) {
    emit('add-selection', item)
    searchTerm.value = ''
  }
}
</script>

<template>
  <div class="search-box">
    <div class="search-box-input">
      <input type="text" placeholder="Search..." v-model="searchTerm" />
      <ul class="suggestions-list">
        <li
          v-for="item in filteredSelections.slice(0, 5)"
          :key="item.id"
          @mousedown="handleSelectItem(item.id)"
        >
          {{ item.value }}
        </li>
      </ul>
    </div>
    <button @click="handleAddItem">add</button>
  </div>
</template>

<style scoped>
.search-box {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: baseline;
  gap: 0.5rem;
}

.search-box-input {
  position: relative;
}

.search-box input {
  width: 100%;
}

.search-box ul {
  list-style: none;
  padding: 0.25rem;
  margin: 0;
  background-color: white;
  border: 1px solid #ccc;
  width: 100%;
}

.search-box li {
  padding-left: 0.25rem;
  padding-right: 0.25rem;
}

@media (hover = hover) {
  .search-box li:hover {
    background-color: green;
    color: white;
  }
}

.search-box .suggestions-list {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
}

.search-box input:focus ~ .suggestions-list {
  display: block;
}
</style>
