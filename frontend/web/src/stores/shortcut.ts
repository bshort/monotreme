import { isEqual } from "lodash-es";
import { create } from "zustand";
import { combine } from "zustand/middleware";
import { shortcutServiceClient } from "@/grpcweb";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";

interface State {
  shortcutMapById: Record<number, Shortcut>;
  shortcutMapByUuid: Record<string, Shortcut>;
  hasMore: boolean;
  totalCount: number;
  isLoadingMore: boolean;
}

const getDefaultState = (): State => {
  return {
    shortcutMapById: {},
    shortcutMapByUuid: {},
    hasMore: false,
    totalCount: 0,
    isLoadingMore: false,
  };
};

const useShortcutStore = create(
  combine(getDefaultState(), (set, get) => ({
    fetchShortcutList: async (reset = true) => {
      if (reset) {
        set({ shortcutMapById: {}, shortcutMapByUuid: {} });
      }

      const { shortcuts, hasMore, totalCount } = await shortcutServiceClient.listShortcuts({
        limit: 50,
        offset: 0,
      });

      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      shortcuts.forEach((shortcut) => {
        shortcutMapById[shortcut.id] = shortcut;
        shortcutMapByUuid[shortcut.uuid] = shortcut;
      });
      set({
        shortcutMapById: shortcutMapById,
        shortcutMapByUuid: shortcutMapByUuid,
        hasMore,
        totalCount,
      });
      return shortcuts;
    },
    loadMoreShortcuts: async () => {
      const { isLoadingMore, hasMore } = get();
      if (isLoadingMore || !hasMore) {
        return;
      }

      set({ isLoadingMore: true });

      try {
        const currentShortcuts = get().getShortcutList();
        const offset = currentShortcuts.length;

        const { shortcuts, hasMore: newHasMore, totalCount } = await shortcutServiceClient.listShortcuts({
          limit: 50,
          offset,
        });

        const shortcutMapById = get().shortcutMapById;
        const shortcutMapByUuid = get().shortcutMapByUuid;
        shortcuts.forEach((shortcut) => {
          shortcutMapById[shortcut.id] = shortcut;
          shortcutMapByUuid[shortcut.uuid] = shortcut;
        });

        set({
          shortcutMapById: shortcutMapById,
          shortcutMapByUuid: shortcutMapByUuid,
          hasMore: newHasMore,
          totalCount,
          isLoadingMore: false,
        });
      } catch (error) {
        set({ isLoadingMore: false });
        throw error;
      }
    },
    fetchShortcutByName: async (name: string) => {
      const shortcut = await shortcutServiceClient.getShortcutByName({
        name,
      });
      return shortcut;
    },
    getOrFetchShortcutById: async (id: number) => {
      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      if (shortcutMapById[id]) {
        return shortcutMapById[id] as Shortcut;
      }

      const shortcut = await shortcutServiceClient.getShortcut({
        id,
      });
      shortcutMapById[id] = shortcut;
      shortcutMapByUuid[shortcut.uuid] = shortcut;
      set({ shortcutMapById: shortcutMapById, shortcutMapByUuid: shortcutMapByUuid });
      return shortcut;
    },
    getShortcutById: (id: number) => {
      const shortcutMap = get().shortcutMapById;
      return shortcutMap[id] || unknownShortcut;
    },
    getOrFetchShortcutByUuid: async (uuid: string) => {
      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      if (shortcutMapByUuid[uuid]) {
        return shortcutMapByUuid[uuid] as Shortcut;
      }

      // Since there's no getShortcutByUuid in the API, we need to fetch all shortcuts first
      await get().fetchShortcutList();
      return shortcutMapByUuid[uuid] || unknownShortcut;
    },
    getShortcutByUuid: (uuid: string) => {
      const shortcutMap = get().shortcutMapByUuid;
      return shortcutMap[uuid] || unknownShortcut;
    },
    getShortcutList: () => {
      return Object.values(get().shortcutMapById);
    },
    createShortcut: async (shortcut: Shortcut) => {
      const createdShortcut = await shortcutServiceClient.createShortcut({
        shortcut: shortcut,
      });
      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      shortcutMapById[createdShortcut.id] = createdShortcut;
      shortcutMapByUuid[createdShortcut.uuid] = createdShortcut;
      set({ shortcutMapById: shortcutMapById, shortcutMapByUuid: shortcutMapByUuid });
      return createdShortcut;
    },
    updateShortcut: async (shortcut: Partial<Shortcut>, updateMask: string[]) => {
      const updatedShortcut = await shortcutServiceClient.updateShortcut({
        shortcut: shortcut,
        updateMask,
      });
      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      shortcutMapById[updatedShortcut.id] = updatedShortcut;
      shortcutMapByUuid[updatedShortcut.uuid] = updatedShortcut;
      set({ shortcutMapById: shortcutMapById, shortcutMapByUuid: shortcutMapByUuid });
      return updatedShortcut;
    },
    deleteShortcut: async (id: number) => {
      const shortcutMapById = get().shortcutMapById;
      const shortcutMapByUuid = get().shortcutMapByUuid;
      const shortcutToDelete = shortcutMapById[id];

      await shortcutServiceClient.deleteShortcut({
        id,
      });

      delete shortcutMapById[id];
      if (shortcutToDelete) {
        delete shortcutMapByUuid[shortcutToDelete.uuid];
      }
      set({ shortcutMapById: shortcutMapById, shortcutMapByUuid: shortcutMapByUuid });
    },
  })),
);

const unknownShortcut: Shortcut = Shortcut.fromPartial({
  id: -1,
  name: "Unknown",
});

export const getShortcutUpdateMask = (shortcut: Shortcut, updatingShortcut: Shortcut) => {
  const updateMask: string[] = [];
  if (!isEqual(shortcut.name, updatingShortcut.name)) {
    updateMask.push("name");
  }
  if (!isEqual(shortcut.link, updatingShortcut.link)) {
    updateMask.push("link");
  }
  if (!isEqual(shortcut.title, updatingShortcut.title)) {
    updateMask.push("title");
  }
  if (!isEqual(shortcut.description, updatingShortcut.description)) {
    updateMask.push("description");
  }
  if (!isEqual(shortcut.tags, updatingShortcut.tags)) {
    updateMask.push("tags");
  }
  if (!isEqual(shortcut.visibility, updatingShortcut.visibility)) {
    updateMask.push("visibility");
  }
  if (!isEqual(shortcut.ogMetadata, updatingShortcut.ogMetadata)) {
    updateMask.push("og_metadata");
  }
  if (!isEqual(shortcut.isFavorite, updatingShortcut.isFavorite)) {
    updateMask.push("is_favorite");
  }
  return updateMask;
};

export default useShortcutStore;
