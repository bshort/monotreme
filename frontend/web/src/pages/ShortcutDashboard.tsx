import { Button, Input } from "@mui/joy";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router-dom";
import useLocalStorage from "react-use/lib/useLocalStorage";
import CreateShortcutDrawer from "@/components/CreateShortcutDrawer";
import FilterView from "@/components/FilterView";
import Icon from "@/components/Icon";
import ShortcutsContainer from "@/components/ShortcutsContainer";
import ShortcutsNavigator from "@/components/ShortcutsNavigator";
import StandaloneViewControls from "@/components/StandaloneViewControls";
import useLoading from "@/hooks/useLoading";
import { useShortcutStore, useUserStore, useViewStore } from "@/stores";
import { getFilteredShortcutList, getOrderedShortcutList } from "@/stores/view";

interface State {
  showCreateShortcutDrawer: boolean;
}

const ShortcutDashboard: React.FC = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const [, setLastVisited] = useLocalStorage<string>("lastVisited", "/shortcuts");
  const loadingState = useLoading();
  const currentUser = useUserStore().getCurrentUser();
  const shortcutStore = useShortcutStore();
  const viewStore = useViewStore();
  const shortcutList = shortcutStore.getShortcutList();
  const [state, setState] = useState<State>({
    showCreateShortcutDrawer: false,
  });

  // Get tags from URL querystring
  const urlTagsParam = searchParams.get('tags');
  const urlTags = urlTagsParam ? urlTagsParam.split(',').map(tag => tag.trim()) : [];

  // Merge URL tags with current filter
  const filter = {
    ...viewStore.filter,
    urlTags: urlTags.length > 0 ? urlTags : undefined
  };

  const filteredShortcutList = getFilteredShortcutList(shortcutList, filter, currentUser);
  const orderedShortcutList = getOrderedShortcutList(filteredShortcutList, viewStore.order);

  useEffect(() => {
    setLastVisited("/shortcuts");
    Promise.all([shortcutStore.fetchShortcutList()]).finally(() => {
      loadingState.setFinish();
    });
  }, []);

  const setShowCreateShortcutDrawer = (show: boolean) => {
    setState({
      ...state,
      showCreateShortcutDrawer: show,
    });
  };

  const handleReload = () => {
    loadingState.setLoading();
    Promise.all([shortcutStore.fetchShortcutList()]).finally(() => {
      loadingState.setFinish();
    });
  };

  return (
    <>
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
        <ShortcutsNavigator />
        <div className="w-full flex flex-row justify-between items-center mb-4">
          <div className="flex flex-row justify-start items-center">
            <Input
              className="w-32 mr-2"
              type="text"
              size="sm"
              placeholder={t("common.search")}
              startDecorator={<Icon.Search className="w-4 h-auto" />}
              value={filter.search}
              onChange={(e) => viewStore.setFilter({ search: e.target.value })}
            />
          </div>
          <div className="flex flex-row justify-end items-center gap-2">
            <Button className="hover:shadow" variant="plain" size="sm" onClick={handleReload} disabled={loadingState.isLoading}>
              <Icon.RotateCcw className="w-4 h-auto" />
              <span className="ml-0.5">{t("common.reload")}</span>
            </Button>
            <Button className="hover:shadow" variant="soft" size="sm" onClick={() => setShowCreateShortcutDrawer(true)}>
              <Icon.Plus className="w-5 h-auto" />
              <span className="ml-0.5">{t("common.create")}</span>
            </Button>
          </div>
        </div>

        {/* Standalone View Controls */}
        <div className="w-full mb-4 flex flex-row justify-start items-center overflow-x-auto">
          <StandaloneViewControls />
        </div>

        <FilterView />
        {loadingState.isLoading ? (
          <div className="py-12 w-full flex flex-row justify-center items-center opacity-80 dark:text-gray-500">
            <Icon.Loader className="mr-2 w-5 h-auto animate-spin" />
            {t("common.loading")}
          </div>
        ) : orderedShortcutList.length === 0 ? (
          <div className="py-16 w-full flex flex-col justify-center items-center text-gray-400">
            <Icon.PackageOpen size={64} strokeWidth={1} />
            <p className="mt-2">No shortcuts found.</p>
            <a
              className="text-blue-600 border-t dark:border-t-zinc-600 text-sm hover:underline flex flex-row justify-center items-center mt-4 pt-2"
              href="https://github.com/bshort/monotreme/blob/main/docs/getting-started/shortcuts.md"
              target="_blank"
            >
              <span>Learn more about shortcuts.</span>
              <Icon.ExternalLink className="ml-1 w-4 h-auto inline" />
            </a>
          </div>
        ) : (
          <ShortcutsContainer shortcutList={orderedShortcutList} />
        )}
      </div>

      {state.showCreateShortcutDrawer && (
        <CreateShortcutDrawer onClose={() => setShowCreateShortcutDrawer(false)} onConfirm={() => setShowCreateShortcutDrawer(false)} />
      )}
    </>
  );
};

export default ShortcutDashboard;
