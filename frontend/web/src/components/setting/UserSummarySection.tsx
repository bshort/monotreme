import { Card, Typography, Chip } from "@mui/joy";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useCollectionStore, useShortcutStore, useUserStore } from "@/stores";
import { getAllUniqueTags } from "@/stores/view";
import Icon from "@/components/Icon";

const UserSummarySection = () => {
  const userStore = useUserStore();
  const shortcutStore = useShortcutStore();
  const collectionStore = useCollectionStore();
  const currentUser = userStore.getCurrentUser();

  const [loading, setLoading] = useState(true);
  const [userList, setUserList] = useState<any[]>([]);

  // Get current data
  const shortcuts = shortcutStore.getShortcutList();
  const collections = collectionStore.getCollectionList();

  // Calculate user's shortcuts and stats
  const userShortcuts = shortcuts.filter(shortcut => shortcut.creatorId === currentUser?.id);
  const userCollections = collections.filter(collection => collection.creatorId === currentUser?.id);
  const userTags = getAllUniqueTags(userShortcuts);
  const totalClicks = userShortcuts.reduce((total, shortcut) => total + shortcut.viewCount, 0);

  // Get list of all users (to show users the current user may have invited)
  useEffect(() => {
    const loadData = async () => {
      try {
        // Load shortcuts and collections
        await Promise.all([
          shortcutStore.fetchShortcutList(),
          collectionStore.fetchCollectionList(),
          userStore.fetchUserList(),
        ]);

        // Get all users except current user
        const allUsers = Object.values(userStore.userMapById).filter(user => user.id !== currentUser?.id);
        setUserList(allUsers);
      } catch (error) {
        console.error("Failed to load user summary data:", error);
      } finally {
        setLoading(false);
      }
    };

    // Only load data once when component mounts
    if (currentUser?.id) {
      loadData();
    }
  }, []); // Empty dependency array to run only once on mount

  if (loading) {
    return (
      <Card className="w-full p-6">
        <Typography level="title-md" className="mb-3">
          Summary
        </Typography>
        <div className="flex items-center justify-center py-8">
          <Icon.Loader className="w-5 h-5 animate-spin mr-2" />
          <span className="text-sm text-gray-500">Loading summary...</span>
        </div>
      </Card>
    );
  }

  return (
    <Card className="w-full p-6">
      <Typography level="title-md" className="mb-4">
        Summary
      </Typography>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Left Column */}
        <div className="space-y-4">
          {/* Shortcuts Count */}
          <Link to="/shortcuts" className="block">
            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors cursor-pointer">
              <div className="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center">
                <Icon.Link className="w-5 h-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
                  {userShortcuts.length} Shortcuts
                </Typography>
                <Typography level="body-xs" className="text-gray-500">
                  Links you've created
                </Typography>
              </div>
            </div>
          </Link>

          {/* Collections Count */}
          <Link to="/collections" className="block">
            <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors cursor-pointer">
              <div className="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center">
                <Icon.FolderOpen className="w-5 h-5 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
                  {userCollections.length} Collections
                </Typography>
                <Typography level="body-xs" className="text-gray-500">
                  Organized groups
                </Typography>
              </div>
            </div>
          </Link>

          {/* Total Clicks */}
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-full flex items-center justify-center">
              <Icon.MousePointer className="w-5 h-5 text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
                {totalClicks.toLocaleString()} Total Clicks
              </Typography>
              <Typography level="body-xs" className="text-gray-500">
                Across all your links
              </Typography>
            </div>
          </div>
        </div>

        {/* Right Column */}
        <div className="space-y-4">
          {/* Tags */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <Icon.Tag className="w-4 h-4 text-gray-600 dark:text-gray-400" />
              <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
                Your Tags ({userTags.length})
              </Typography>
            </div>
            <div className="flex flex-wrap gap-1 max-h-24 overflow-y-auto">
              {userTags.length > 0 ? (
                userTags.slice(0, 10).map((tag) => (
                  <Link key={tag} to={`/shortcuts?tags=${encodeURIComponent(tag)}`}>
                    <Chip size="sm" variant="soft" color="primary" className="hover:bg-blue-200 dark:hover:bg-blue-800 cursor-pointer">
                      #{tag}
                    </Chip>
                  </Link>
                ))
              ) : (
                <Typography level="body-xs" className="text-gray-400 italic">
                  No tags yet
                </Typography>
              )}
              {userTags.length > 10 && (
                <Chip size="sm" variant="outlined" color="neutral">
                  +{userTags.length - 10} more
                </Chip>
              )}
            </div>
          </div>

          {/* Other Users */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <Icon.Users className="w-4 h-4 text-gray-600 dark:text-gray-400" />
              <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
                Other Users ({userList.length})
              </Typography>
            </div>
            <div className="space-y-1 max-h-20 overflow-y-auto">
              {userList.length > 0 ? (
                userList.slice(0, 5).map((user) => (
                  <div key={user.id} className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    <Typography level="body-xs" className="text-gray-600 dark:text-gray-400 truncate">
                      {user.nickname} ({user.email})
                    </Typography>
                  </div>
                ))
              ) : (
                <Typography level="body-xs" className="text-gray-400 italic">
                  No other users
                </Typography>
              )}
              {userList.length > 5 && (
                <Typography level="body-xs" className="text-gray-500">
                  +{userList.length - 5} more users
                </Typography>
              )}
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default UserSummarySection;