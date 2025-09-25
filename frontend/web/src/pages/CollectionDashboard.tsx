import { Button, Input, Select, Option, Tooltip } from "@mui/joy";
import copy from "copy-to-clipboard";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import useLocalStorage from "react-use/lib/useLocalStorage";
import CollectionView from "@/components/CollectionView";
import CreateCollectionDrawer from "@/components/CreateCollectionDrawer";
import FilterView from "@/components/FilterView";
import Icon from "@/components/Icon";
import { userServiceClient } from "@/grpcweb";
import useLoading from "@/hooks/useLoading";
import { useShortcutStore, useCollectionStore, useUserStore } from "@/stores";
import { Collection } from "@/types/proto/api/v1/collection_service";

interface State {
  showCreateCollectionDrawer: boolean;
}

type SortOption = "date" | "name" | "shortcuts";

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
  const sortCollections = (collections: Collection[], sortBy: SortOption): Collection[] => {
    return [...collections].sort((a, b) => {
      switch (sortBy) {
        case "name":
          return a.title.toLowerCase().localeCompare(b.title.toLowerCase());
        case "shortcuts":
          return b.shortcutIds.length - a.shortcutIds.length;
        case "date":
        default:
          return b.createdTs - a.createdTs;
      }
    });
  };

  const filteredAndSortedCollections = sortCollections(
    collectionStore.getCollectionList().filter((collection) => {
      return (
        collection.name.toLowerCase().includes(search.toLowerCase()) ||
        collection.title.toLowerCase().includes(search.toLowerCase()) ||
        collection.description.toLowerCase().includes(search.toLowerCase())
      );
    }),
    sortBy
  );

  useEffect(() => {
    setLastVisited("/collections");
    Promise.all([shortcutStore.fetchShortcutList(), collectionStore.fetchCollectionList()]).finally(() => {
      loadingState.setFinish();
    });
  }, []);

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
            </Select>
          </div>
          <div className="flex flex-row justify-start items-center gap-2">
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
        ) : (
          <div className="w-full flex flex-col justify-start items-start gap-3">
            {filteredAndSortedCollections.map((collection) => {
              return <CollectionView key={collection.id} collection={collection} />;
            })}
          </div>
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
