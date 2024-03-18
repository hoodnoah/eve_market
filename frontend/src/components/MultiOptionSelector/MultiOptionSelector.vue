<script setup lang="ts">
// types
import type { IValueItem } from '@/types/ValueItem'

// components
import SelectorOption from './SelectorOption.vue'
import SearchBox from './SearchBox.vue'

const props = defineProps({
  header: {
    type: String,
    required: true
  }
})

const optionsModel = defineModel<IValueItem[]>({ required: true })

function clearSelection() {
  optionsModel.value = []
}

function handleRemove(event: Event) {
  const id = Number(event)
  optionsModel.value = optionsModel.value.filter((option) => option.id !== Number(id))
}

function handleAddItem(event: Event) {
  const item = event as unknown as IValueItem
  optionsModel.value.push(item)
}
</script>

<template>
  <div class="multi-option-selector">
    <div class="header">
      <h3>{{ props.header }}</h3>
      <button @click="clearSelection">clear</button>
    </div>
    <div class="controls">
      <SearchBox @add-selection="handleAddItem"></SearchBox>
    </div>
    <SelectorOption
      v-for="option in optionsModel"
      :key="option.id"
      :selection="option"
      @remove-selection="handleRemove"
    ></SelectorOption>
  </div>
</template>

<style scoped>
.multi-option-selector {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;

  padding: 0.75rem;

  border: 1px solid #ccc;
  border-radius: 0.25rem;
}

.multi-option-selector .header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: baseline;
}

.multi-option-selector .header h3 {
  flex-grow: 1;
  flex-shrink: 1;
}

.multi-option-selector .header button {
  flex-grow: 0;
  flex-shrink: 0;
}

.multi-option-selector .controls {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: baseline;
  gap: 0.25rem;

  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.multi-option-selector .controls .search-box {
  flex-grow: 1;
  flex-shrink: 1;
}

.multi-option-selector .controls button {
  flex-grow: 0;
  flex-shrink: 0;
}
</style>
