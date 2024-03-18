<script setup lang="ts">
import { ref } from 'vue'

// types
import type { RegionSelection, ItemSelection } from './types/ValueItem'

// components
import MultiOptionSelector from './components/MultiOptionSelector/MultiOptionSelector.vue'

// state objects
const selectedItems = ref<ItemSelection[]>([])
const selectedRegions = ref<RegionSelection[]>([])

selectedItems.value = [
  { id: 1, value: 'Option 1' },
  { id: 2, value: 'Option 2' },
  { id: 3, value: 'Option 3' },
  { id: 4, value: 'Option 4' },
  { id: 5, value: 'Option 5' }
]

selectedRegions.value = [
  { id: 1, value: 'Region 1' },
  { id: 2, value: 'Region 2' },
  { id: 3, value: 'Region 3' },
  { id: 4, value: 'Region 4' },
  { id: 5, value: 'Region 5' }
]
</script>

<template>
  <header>
    <div class="wrapper">
      <h1>eve market</h1>
    </div>
  </header>

  <main>
    <div class="selections-wrapper">
      <MultiOptionSelector header="items" :modelValue="selectedItems" />
      <MultiOptionSelector header="regions" :modelValue="selectedRegions" />
    </div>
  </main>
</template>

<style scoped>
.selections-wrapper {
  display: flex;
  flex-direction: row;
  gap: 0.75rem;
}
</style>
