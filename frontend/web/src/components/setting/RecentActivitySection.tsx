import { Card, Typography, Chip, Divider, Avatar } from "@mui/joy";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useUserStore, useShortcutStore, useCollectionStore } from "@/stores";
import { userServiceClient, shortcutServiceClient, collectionServiceClient } from "@/grpcweb";
import { User } from "@/types/proto/api/v1/user_service";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { Collection } from "@/types/proto/api/v1/collection_service";
import Icon from "@/components/Icon";

interface RecentActivity {
  id: string;
  type: 'user' | 'shortcut' | 'collection' | 'click';
  data: any;
  timestamp: Date;
}

const RecentActivitySection = () => {
  const userStore = useUserStore();
  const [loading, setLoading] = useState(true);
  const [recentUsers, setRecentUsers] = useState<User[]>([]);
  const [recentShortcuts, setRecentShortcuts] = useState<Shortcut[]>([]);
  const [recentCollections, setRecentCollections] = useState<Collection[]>([]);
  const [recentClicks, setRecentClicks] = useState<any[]>([]);

  useEffect(() => {
    const loadRecentActivity = async () => {
      try {
        setLoading(true);

        // Fetch all data
        const [usersResponse, shortcutsResponse, collectionsResponse] = await Promise.all([
          userServiceClient.listUsers({}),
          shortcutServiceClient.listShortcuts({}),
          collectionServiceClient.listCollections({})
        ]);

        // Sort and limit to 5 most recent items
        const sortedUsers = usersResponse.users
          .sort((a, b) => new Date(b.createdTime).getTime() - new Date(a.createdTime).getTime())
          .slice(0, 5);

        const sortedShortcuts = shortcutsResponse.shortcuts
          .sort((a, b) => new Date(b.createdTime).getTime() - new Date(a.createdTime).getTime())
          .slice(0, 5);

        const sortedCollections = collectionsResponse.collections
          .sort((a, b) => new Date(b.createdTime).getTime() - new Date(a.createdTime).getTime())
          .slice(0, 5);

        // For clicks, we'll use the shortcuts with highest view counts as proxy for recent activity
        // This is a simplification since we don't have direct access to click/activity API
        const recentClickActivity = shortcutsResponse.shortcuts
          .filter(s => s.viewCount > 0)
          .sort((a, b) => b.viewCount - a.viewCount)
          .slice(0, 5)
          .map(shortcut => ({
            id: `click-${shortcut.id}`,
            shortcut,
            timestamp: new Date() // Placeholder since we don't have actual click timestamps
          }));

        setRecentUsers(sortedUsers);
        setRecentShortcuts(sortedShortcuts);
        setRecentCollections(sortedCollections);
        setRecentClicks(recentClickActivity);

      } catch (error) {
        console.error("Failed to load recent activity:", error);
      } finally {
        setLoading(false);
      }
    };

    loadRecentActivity();
  }, []);

  const formatTimeAgo = (date: Date) => {
    const now = new Date();
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));

    if (diffInHours < 1) return "Just now";
    if (diffInHours < 24) return `${diffInHours}h ago`;
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7) return `${diffInDays}d ago`;
    return date.toLocaleDateString();
  };

  const getInitials = (name: string) => {
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  if (loading) {
    return (
      <Card className="w-full p-6">
        <Typography level="title-md" className="mb-3">
          Recent Activity
        </Typography>
        <div className="flex items-center justify-center py-8">
          <Icon.Loader className="w-5 h-5 animate-spin mr-2" />
          <span className="text-sm text-gray-500">Loading recent activity...</span>
        </div>
      </Card>
    );
  }

  return (
    <Card className="w-full p-6">
      <Typography level="title-md" className="mb-4">
        Recent Activity
      </Typography>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Recent Users */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.Users className="w-4 h-4 text-blue-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Users
            </Typography>
          </div>
          <div className="space-y-2">
            {recentUsers.length > 0 ? (
              recentUsers.map((user) => (
                <div key={user.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <Avatar size="sm" className="bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                    {getInitials(user.nickname || user.email)}
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <Typography level="body-sm" className="truncate font-medium">
                      {user.nickname || user.email}
                    </Typography>
                    <Typography level="body-xs" className="text-gray-500 truncate">
                      {formatTimeAgo(new Date(user.createdTime))}
                    </Typography>
                  </div>
                </div>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent users
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Shortcuts */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.Link className="w-4 h-4 text-green-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Shortcuts
            </Typography>
          </div>
          <div className="space-y-2">
            {recentShortcuts.length > 0 ? (
              recentShortcuts.map((shortcut) => (
                <Link key={shortcut.id} to={`/shortcut/${shortcut.id}`} className="block">
                  <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <div className="w-8 h-8 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Icon.Link className="w-4 h-4 text-green-600 dark:text-green-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <Typography level="body-sm" className="truncate font-medium">
                        {shortcut.title || shortcut.name}
                      </Typography>
                      <Typography level="body-xs" className="text-gray-500 truncate">
                        {formatTimeAgo(new Date(shortcut.createdTime))}
                      </Typography>
                    </div>
                  </div>
                </Link>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent shortcuts
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Collections */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.FolderOpen className="w-4 h-4 text-purple-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Collections
            </Typography>
          </div>
          <div className="space-y-2">
            {recentCollections.length > 0 ? (
              recentCollections.map((collection) => (
                <div key={collection.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Icon.FolderOpen className="w-4 h-4 text-purple-600 dark:text-purple-400" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <Typography level="body-sm" className="truncate font-medium">
                      {collection.title || collection.name}
                    </Typography>
                    <Typography level="body-xs" className="text-gray-500 truncate">
                      {formatTimeAgo(new Date(collection.createdTime))}
                    </Typography>
                  </div>
                </div>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent collections
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Clicks */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.MousePointer className="w-4 h-4 text-orange-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Most Clicked
            </Typography>
          </div>
          <div className="space-y-2">
            {recentClicks.length > 0 ? (
              recentClicks.map((click) => (
                <Link key={click.id} to={`/shortcut/${click.shortcut.id}`} className="block">
                  <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <div className="w-8 h-8 bg-orange-100 dark:bg-orange-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Icon.MousePointer className="w-4 h-4 text-orange-600 dark:text-orange-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <Typography level="body-sm" className="truncate font-medium">
                        {click.shortcut.title || click.shortcut.name}
                      </Typography>
                      <Typography level="body-xs" className="text-gray-500 truncate">
                        {click.shortcut.viewCount} clicks
                      </Typography>
                    </div>
                  </div>
                </Link>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No clicks yet
              </Typography>
            )}
          </div>
        </div>
      </div>
    </Card>
  );
};

export default RecentActivitySection;