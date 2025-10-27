import { Button, Input } from "@mui/joy";
import { useEffect, useRef, useState } from "react";
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
  tagsExpanded: boolean;
  visibleTagsCount: number;
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
    tagsExpanded: false,
    visibleTagsCount: 0,
  });
  const tagsContainerRef = useRef<HTMLDivElement>(null);

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

  // Get tags sorted by popularity
  const sortedTags = getSortedTagsByPopularity(shortcutList);

  useEffect(() => {
    setLastVisited("/shortcuts");
    Promise.all([shortcutStore.fetchShortcutList()]).finally(() => {
      loadingState.setFinish();
    });
  }, []);

  // Calculate visible tags count when collapsed
  useEffect(() => {
    if (!state.tagsExpanded && tagsContainerRef.current && sortedTags.length > 0) {
      const container = tagsContainerRef.current;
      const children = Array.from(container.children) as HTMLElement[];

      let visibleCount = 0;
      const maxHeight = 32; // max-h-8 = 2rem = 32px

      for (const child of children) {
        const rect = child.getBoundingClientRect();
        const containerRect = container.getBoundingClientRect();

        // Check if the tag is within the visible height
        if (rect.top - containerRect.top < maxHeight) {
          visibleCount++;
        } else {
          break;
        }
      }

      if (visibleCount !== state.visibleTagsCount) {
        setState((prevState) => ({ ...prevState, visibleTagsCount: visibleCount }));
      }
    }
  }, [sortedTags, state.tagsExpanded, state.visibleTagsCount]);

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

  const toggleTagsExpanded = () => {
    setState({ ...state, tagsExpanded: !state.tagsExpanded });
  };

  const handleTagClick = (tag: string) => {
    window.location.href = `/shortcuts?tags=${encodeURIComponent(tag)}`;
  };

  return (
    <>
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
        <ShortcutsNavigator />

        {/* Collapsible Tags Section */}
        {sortedTags.length > 0 && (
          <div className="w-full mb-4 border border-gray-200 dark:border-gray-700 rounded-lg">
            <button
              onClick={toggleTagsExpanded}
              className="w-full px-4 py-3 flex items-start justify-between hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              <div className="flex flex-col items-start gap-2 flex-1 overflow-hidden">
                <div className="flex items-center gap-2">
                  <Icon.Tag className="w-4 h-4 text-gray-500 dark:text-gray-400 flex-shrink-0" />
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300 flex-shrink-0">Tags</span>
                </div>
                <div className={`w-full flex flex-wrap gap-2 ${!state.tagsExpanded ? 'max-h-8 overflow-hidden' : ''}`} ref={tagsContainerRef}>
                  {sortedTags.map((tag) => (
                    <span
                      key={tag}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleTagClick(tag);
                      }}
                      className="text-xs px-2.5 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-full hover:bg-blue-100 dark:hover:bg-blue-900/50 cursor-pointer transition-colors font-medium"
                    >
                      #{tag}
                    </span>
                  ))}
                </div>
                {!state.tagsExpanded && state.visibleTagsCount > 0 && state.visibleTagsCount < sortedTags.length && (
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    and {sortedTags.length - state.visibleTagsCount} more...
                  </span>
                )}
              </div>
              <Icon.ChevronDown
                className={`w-4 h-4 text-gray-500 dark:text-gray-400 flex-shrink-0 transition-transform ml-2 ${state.tagsExpanded ? 'rotate-180' : ''}`}
              />
            </button>
          </div>
        )}

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
            <p className="mt-2">No bookmarks found.</p>
            <a
              className="text-blue-600 border-t dark:border-t-zinc-600 text-sm hover:underline flex flex-row justify-center items-center mt-4 pt-2"
              href="https://github.com/bshort/monotreme/blob/main/docs/getting-started/shortcuts.md"
              target="_blank"
            >
              <span>Learn more about bookmarks.</span>
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

// Helper function to sort tags by popularity (most used first)
const getSortedTagsByPopularity = (shortcuts: any[]): string[] => {
  const tagCounts = new Map<string, number>();

  shortcuts.forEach((shortcut) => {
    if (shortcut.tags && Array.isArray(shortcut.tags)) {
      shortcut.tags.forEach((tag: string) => {
        const cleanTag = tag.trim().replace(/,$/, ""); // Remove trailing comma
        if (cleanTag) {
          tagCounts.set(cleanTag, (tagCounts.get(cleanTag) || 0) + 1);
        }
      });
    }
  });

  return Array.from(tagCounts.entries())
    .sort((a, b) => b[1] - a[1]) // Sort by count descending
    .map(([tag]) => tag);
};

export default ShortcutDashboard;
