<script setup lang="ts">
import { type PropType } from 'vue'

const emit = defineEmits(['remove-selection'])

type ValueItem = {
  id: number
  value: string
}

const props = defineProps({
  selection: {
    type: Object as PropType<ValueItem>,
    required: true
  }
})
</script>

<template>
  <div class="selector-option">
    <p>{{ props.selection.value }}</p>
    <button @click="emit('remove-selection', props.selection.id)">remove</button>
  </div>
</template>

<style scoped>
.selector-option {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: baseline;
}
</style>
