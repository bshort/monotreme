import { Button, Input, Select, Option, Tooltip } from "@mui/joy";
import copy from "copy-to-clipboard";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import useLocalStorage from "react-use/lib/useLocalStorage";
import { DndContext, closestCenter, PointerSensor, useSensor, useSensors, DragEndEvent } from "@dnd-kit/core";
import { arrayMove, SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import CollectionView from "@/components/CollectionView";
import CreateCollectionDrawer from "@/components/CreateCollectionDrawer";
import FilterView from "@/components/FilterView";
import Icon from "@/components/Icon";
import SortableCollectionItem from "@/components/SortableCollectionItem";
import { userServiceClient } from "@/grpcweb";
import useLoading from "@/hooks/useLoading";
import { useShortcutStore, useCollectionStore, useUserStore } from "@/stores";
import { Collection } from "@/types/proto/api/v1/collection_service";

interface State {
  showCreateCollectionDrawer: boolean;
}

type SortOption = "date" | "name" | "shortcuts" | "custom";

const CollectionDashboard: React.FC = () => {
  const { t } = useTranslation();
  const [, setLastVisited] = useLocalStorage<string>("lastVisited", "/shortcuts");
  const loadingState = useLoading();
  const shortcutStore = useShortcutStore();
  const collectionStore = useCollectionStore();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const [state, setState] = useState<State>({
    showCreateCollectionDrawer: false,
  });
  const [search, setSearch] = useState<string>("");
  const [sortBy, setSortBy] = useState<SortOption>("date");
  const [collectionList, setCollectionList] = useState<Collection[]>([]);
  const sortCollections = (collections: Collection[], sortBy: SortOption): Collection[] => {
    return [...collections].sort((a, b) => {
      switch (sortBy) {
        case "name":
          return a.title.toLowerCase().localeCompare(b.title.toLowerCase());
        case "shortcuts":
          return b.shortcutIds.length - a.shortcutIds.length;
        case "custom":
          return a.userOrder - b.userOrder;
        case "date":
        default:
          return b.createdTs - a.createdTs;
      }
    });
  };

  const filteredAndSortedCollections = sortCollections(
    collectionList.filter((collection) => {
      const searchLower = search.toLowerCase();

      // Search collection properties
      const matchesCollection =
        collection.name.toLowerCase().includes(searchLower) ||
        collection.title.toLowerCase().includes(searchLower) ||
        collection.description.toLowerCase().includes(searchLower);

      if (matchesCollection) {
        return true;
      }

      // Search shortcuts within the collection
      const matchesShortcut = collection.shortcutIds.some((shortcutId) => {
        const shortcut = shortcutStore.shortcutMapById[shortcutId];
        if (!shortcut) {
          return false;
        }

        return (
          shortcut.name.toLowerCase().includes(searchLower) ||
          shortcut.title.toLowerCase().includes(searchLower) ||
          shortcut.description.toLowerCase().includes(searchLower) ||
          shortcut.link.toLowerCase().includes(searchLower) ||
          shortcut.tags.some((tag) => tag.toLowerCase().includes(searchLower))
        );
      });

      return matchesShortcut;
    }),
    sortBy
  );

  useEffect(() => {
    setLastVisited("/collections");
    Promise.all([shortcutStore.fetchShortcutList(), collectionStore.fetchCollectionList()]).finally(() => {
      loadingState.setFinish();
    });
  }, []);

  useEffect(() => {
    setCollectionList(collectionStore.getCollectionList());
  }, [collectionStore.collectionMapById]);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const isCustomOrder = sortBy === "custom";

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    if (!over || active.id === over.id) {
      return;
    }

    const oldIndex = filteredAndSortedCollections.findIndex((c) => c.id === active.id);
    const newIndex = filteredAndSortedCollections.findIndex((c) => c.id === over.id);

    const newOrderedList = arrayMove(filteredAndSortedCollections, oldIndex, newIndex);

    // Get the full collection list to properly update all userOrder values
    const fullCollectionList = collectionList;

    // Create a map of id to new userOrder based on the reordered visible list
    const visibleOrderMap = new Map(newOrderedList.map((c, idx) => [c.id, idx]));

    // Separate visible and non-visible collections
    const visibleCollections = newOrderedList;
    const nonVisibleCollections = fullCollectionList.filter((c) => !visibleOrderMap.has(c.id));

    // Assign sequential userOrder values to ALL collections
    // Visible collections get 0, 1, 2, ... based on their new order
    // Non-visible collections get sequential values starting after the visible ones
    const updatedCollections = [
      ...visibleCollections.map((c, idx) => {
        const updated = { ...c };
        updated.userOrder = idx;
        return updated;
      }),
      ...nonVisibleCollections.map((c, idx) => {
        const updated = { ...c };
        updated.userOrder = visibleCollections.length + idx;
        return updated;
      }),
    ];

    // Optimistically update the local state
    setCollectionList(updatedCollections);

    try {
      // Only update collections where userOrder actually changed
      const collectionsToUpdate = updatedCollections.filter((updated) => {
        const original = fullCollectionList.find((c) => c.id === updated.id);
        return !original || original.userOrder !== updated.userOrder;
      });

      // Update sequentially to avoid race conditions
      for (const collection of collectionsToUpdate) {
        await collectionStore.updateCollection(
          { id: collection.id, userOrder: collection.userOrder },
          ["user_order"]
        );
      }

      await collectionStore.fetchCollectionList();
    } catch (error) {
      console.error("Failed to update collection order:", error);
      toast.error("Failed to update collection order");
      setCollectionList(collectionStore.getCollectionList());
    }
  };

  const setShowCreateCollectionDrawer = (show: boolean) => {
    setState({
      ...state,
      showCreateCollectionDrawer: show,
    });
  };

  const handleCopyCollectionsRSSLink = async () => {
    try {
      // Generate a new access token for RSS feed
      const { accessToken } = await userServiceClient.createUserAccessToken({
        id: currentUser.id,
        description: "RSS Feed",
        expiresAt: undefined, // Never expires
      });

      const rssUrl = `${window.location.origin}/rss/collections.xml?token=${accessToken}`;
      copy(rssUrl);
      toast.success("RSS feed URL with access token copied to clipboard!");
    } catch (error: any) {
      console.error("Failed to create RSS access token:", error);
      toast.error("Failed to generate RSS feed URL. Please try again.");
    }
  };

  const handleReload = () => {
    loadingState.setLoading();
    Promise.all([shortcutStore.fetchShortcutList(), collectionStore.fetchCollectionList()]).finally(() => {
      loadingState.setFinish();
    });
  };

  return (
    <>
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
        <div className="w-full flex flex-row justify-between items-center mb-4">
          <div className="flex flex-row items-center gap-2">
            <Input
              className="w-32"
              type="text"
              size="sm"
              placeholder={t("common.search")}
              startDecorator={<Icon.Search className="w-4 h-auto" />}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <Select
              size="sm"
              value={sortBy}
              onChange={(_, value) => setSortBy(value as SortOption)}
              startDecorator={<Icon.ArrowUpDown className="w-4 h-auto" />}
            >
              <Option value="date">By Date</Option>
              <Option value="name">By Name</Option>
              <Option value="shortcuts">By Shortcuts</Option>
              <Option value="custom">Custom Order</Option>
            </Select>
          </div>
          <div className="flex flex-row justify-start items-center gap-2">
            <Button className="hover:shadow" variant="plain" size="sm" onClick={handleReload} disabled={loadingState.isLoading}>
              <Icon.RotateCcw className="w-4 h-auto" />
              <span className="ml-0.5">{t("common.reload")}</span>
            </Button>
            <Tooltip title="Copy RSS feed URL for all collections" placement="top" arrow>
              <Button className="hover:shadow" variant="plain" size="sm" onClick={() => handleCopyCollectionsRSSLink()}>
                <Icon.Rss className="w-4 h-auto" />
              </Button>
            </Tooltip>
            <Button className="hover:shadow" variant="soft" size="sm" onClick={() => setShowCreateCollectionDrawer(true)}>
              <Icon.Plus className="w-5 h-auto" />
              <span className="ml-0.5">{t("common.create")}</span>
            </Button>
          </div>
        </div>
        <FilterView />
        {loadingState.isLoading ? (
          <div className="py-12 w-full flex flex-row justify-center items-center opacity-80 dark:text-gray-500">
            <Icon.Loader className="mr-2 w-5 h-auto animate-spin" />
            {t("common.loading")}
          </div>
        ) : filteredAndSortedCollections.length === 0 ? (
          <div className="py-16 w-full flex flex-col justify-center items-center text-gray-400">
            <Icon.PackageOpen size={64} strokeWidth={1} />
            <p className="mt-2">No collections found.</p>
            <a
              className="text-blue-600 border-t dark:border-t-zinc-600 text-sm hover:underline flex flex-row justify-center items-center mt-4 pt-2"
              href="https://github.com/bshort/monotreme/blob/main/docs/getting-started/collections.md"
              target="_blank"
            >
              <span>Learn more about collections.</span>
              <Icon.ExternalLink className="ml-1 w-4 h-auto inline" />
            </a>
          </div>
        ) : !isCustomOrder ? (
          <div className="w-full flex flex-col justify-start items-start gap-3">
            {filteredAndSortedCollections.map((collection) => {
              return <CollectionView key={collection.id} collection={collection} searchQuery={search} />;
            })}
          </div>
        ) : (
          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
            <SortableContext items={filteredAndSortedCollections.map((c) => c.id)} strategy={verticalListSortingStrategy}>
              <div className="w-full flex flex-col justify-start items-start gap-3">
                {filteredAndSortedCollections.map((collection) => {
                  return <SortableCollectionItem key={collection.id} collection={collection} searchQuery={search} />;
                })}
              </div>
            </SortableContext>
          </DndContext>
        )}
      </div>

      {state.showCreateCollectionDrawer && (
        <CreateCollectionDrawer
          onClose={() => setShowCreateCollectionDrawer(false)}
          onConfirm={() => setShowCreateCollectionDrawer(false)}
        />
      )}
    </>
  );
};

export default CollectionDashboard;
