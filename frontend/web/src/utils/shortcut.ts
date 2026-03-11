import useWorkspaceStore from "@/stores/workspace";

export const getShortcutUrl = (shortcutName: string): string => {
  const workspaceStore = useWorkspaceStore.getState();
  const prefix = workspaceStore.getShortcutPrefix();
  return `/${prefix}/${shortcutName}`;
};

export const getShortcutAbsoluteUrl = (shortcutName: string): string => {
  const workspaceStore = useWorkspaceStore.getState();
  const prefix = workspaceStore.getShortcutPrefix();
  return `${window.location.protocol}//${window.location.host}/${prefix}/${shortcutName}`;
};