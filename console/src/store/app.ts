import { ref } from "vue";
import { defineStore } from "pinia";

export const useAppStore = defineStore("app", () => {
  const sidebarCollapsed = ref(false);

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  }

  function setSidebarCollapsed(value: boolean) {
    sidebarCollapsed.value = value;
  }

  return {
    sidebarCollapsed,
    toggleSidebar,
    setSidebarCollapsed,
  };
});
