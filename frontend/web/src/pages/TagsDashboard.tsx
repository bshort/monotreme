import { Card, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import useLocalStorage from "react-use/lib/useLocalStorage";
import Icon from "@/components/Icon";
import useLoading from "@/hooks/useLoading";
import { useShortcutStore, useUserStore } from "@/stores";
import { getAllUniqueTags } from "@/stores/view";

const TagsDashboard: React.FC = () => {
  const { t } = useTranslation();
  const [, setLastVisited] = useLocalStorage<string>("lastVisited", "/tags");
  const loadingState = useLoading();
  const currentUser = useUserStore().getCurrentUser();
  const shortcutStore = useShortcutStore();
  const shortcutList = shortcutStore.getShortcutList();

  const [tagData, setTagData] = useState<Map<string, Array<{
    id: number;
    name: string;
    title: string;
    link: string;
    viewCount: number;
    ogMetadata?: { imageUrl?: string };
  }>>>(new Map());

  useEffect(() => {
    setLastVisited("/tags");
    Promise.all([shortcutStore.fetchShortcutList()]).finally(() => {
      loadingState.setFinish();
    });
  }, []);

  useEffect(() => {
    // Group shortcuts by tags
    const tagMap = new Map<string, Array<{
      id: number;
      name: string;
      title: string;
      link: string;
      viewCount: number;
      ogMetadata?: { imageUrl?: string };
    }>>();

    shortcutList.forEach((shortcut) => {
      shortcut.tags.forEach((tag) => {
        const cleanTag = tag.trim().replace(/,$/, ''); // Remove trailing comma
        if (cleanTag) {
          if (!tagMap.has(cleanTag)) {
            tagMap.set(cleanTag, []);
          }
          tagMap.get(cleanTag)!.push({
            id: shortcut.id,
            name: shortcut.name,
            title: shortcut.title,
            link: shortcut.link,
            viewCount: shortcut.viewCount,
            ogMetadata: shortcut.ogMetadata
          });
        }
      });
    });

    // Sort tags alphabetically and sort shortcuts within each tag by name
    const sortedTagMap = new Map([...tagMap.entries()].sort());
    sortedTagMap.forEach((shortcuts) => {
      shortcuts.sort((a, b) => a.name.localeCompare(b.name));
    });

    setTagData(sortedTagMap);
  }, [shortcutList]);

  if (loadingState.isLoading) {
    return (
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
        <div className="py-12 w-full flex flex-row justify-center items-center opacity-80 dark:text-gray-500">
          <Icon.Loader className="mr-2 w-5 h-auto animate-spin" />
          {t("common.loading")}
        </div>
      </div>
    );
  }

  if (tagData.size === 0) {
    return (
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
        <div className="py-16 w-full flex flex-col justify-center items-center text-gray-400">
          <Icon.Tag size={64} strokeWidth={1} />
          <p className="mt-2">No tags found.</p>
          <p className="text-sm mt-1">Create some shortcuts with tags to see them organized here.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
      <div className="w-full mb-6">
        <Typography level="h2" className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
          Tags
        </Typography>
        <Typography level="body-md" className="text-gray-600 dark:text-gray-400">
          Browse shortcuts organized by tags
        </Typography>
      </div>

      <div className="w-full space-y-6">
        {Array.from(tagData.entries()).map(([tag, shortcuts]) => (
          <Card key={tag} className="w-full p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <Icon.Tag className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                <Typography level="title-lg" className="text-gray-900 dark:text-gray-100">
                  #{tag}
                </Typography>
                <span className="text-sm text-gray-500 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded-full">
                  {shortcuts.length} shortcut{shortcuts.length !== 1 ? 's' : ''}
                </span>
              </div>
              <Link
                to={`/shortcuts?tags=${encodeURIComponent(tag)}`}
                className="text-blue-600 dark:text-blue-400 hover:underline text-sm flex items-center gap-1"
              >
                View all
                <Icon.ExternalLink className="w-3 h-3" />
              </Link>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {shortcuts.map((shortcut) => (
                <div
                  key={shortcut.id}
                  className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                >
                  {/* Favicon */}
                  <div className="w-8 h-8 flex-shrink-0">
                    {shortcut.ogMetadata?.imageUrl ? (
                      <img
                        className="w-8 h-8 rounded-sm"
                        src={shortcut.ogMetadata.imageUrl}
                        alt={shortcut.title || shortcut.name}
                        onError={(e) => {
                          (e.target as HTMLImageElement).style.display = 'none';
                        }}
                      />
                    ) : (
                      <div className="w-8 h-8 bg-gray-200 dark:bg-gray-700 rounded-sm flex items-center justify-center">
                        <Icon.Globe className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      </div>
                    )}
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                        {shortcut.title || shortcut.name}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 mt-1">
                      <a
                        href={shortcut.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 truncate max-w-48"
                      >
                        {shortcut.link}
                      </a>
                      <div className="flex items-center gap-1 text-xs text-gray-400">
                        <Icon.Eye className="w-3 h-3" />
                        {shortcut.viewCount}
                      </div>
                    </div>
                  </div>

                  {/* Action buttons */}
                  <div className="flex items-center gap-1">
                    <Link
                      to={`/shortcut/${shortcut.id}?edit=true`}
                      className="p-2 text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                      title="Edit shortcut"
                    >
                      <Icon.Edit className="w-4 h-4" />
                    </Link>
                    <a
                      href={shortcut.link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="p-2 text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                      title="Open shortcut"
                    >
                      <Icon.ExternalLink className="w-4 h-4" />
                    </a>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default TagsDashboard;