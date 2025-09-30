import { Button, Card, Input, Select, Switch, Typography, Option, Modal, ModalDialog, ModalClose, Textarea } from "@mui/joy";
import { useState, useEffect } from "react";
import toast from "react-hot-toast";
import Icon from "@/components/Icon";
import { rssImportServiceClient } from "@/grpcweb";
import { RssFeed, Visibility } from "@/types/proto/api/v1/rss_import_service";

interface RssImportSectionProps {}

const RssImportSection: React.FC<RssImportSectionProps> = () => {
  const [rssFeeds, setRssFeeds] = useState<RssFeed[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showImportModal, setShowImportModal] = useState(false);
  const [selectedFeed, setSelectedFeed] = useState<RssFeed | null>(null);
  const [importResult, setImportResult] = useState<any>(null);

  const [createForm, setCreateForm] = useState({
    title: "",
    url: "",
    description: "",
    autoImport: false,
    importFrequencyHours: 24,
    defaultVisibility: Visibility.WORKSPACE,
    defaultTags: "",
    shortcutPrefix: "",
  });

  useEffect(() => {
    loadRssFeeds();
  }, []);

  const loadRssFeeds = async () => {
    try {
      setLoading(true);
      const response = await rssImportServiceClient.listRssFeeds({});
      setRssFeeds(response.rssFeeds || []);
    } catch (error: any) {
      console.error("Failed to load RSS feeds:", error);
      toast.error("Failed to load RSS feeds");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateFeed = async () => {
    try {
      const rssFeed: RssFeed = {
        id: 0,
        creatorId: 0,
        title: createForm.title,
        url: createForm.url,
        description: createForm.description,
        autoImport: createForm.autoImport,
        importFrequencyHours: createForm.importFrequencyHours,
        defaultVisibility: createForm.defaultVisibility,
        defaultTags: createForm.defaultTags.split(",").map(tag => tag.trim()).filter(tag => tag),
        shortcutPrefix: createForm.shortcutPrefix,
        isActive: true,
        lastError: "",
        totalImported: 0,
      };

      await rssImportServiceClient.createRssFeed({ rssFeed });
      toast.success("RSS feed created successfully");
      setShowCreateModal(false);
      setCreateForm({
        title: "",
        url: "",
        description: "",
        autoImport: false,
        importFrequencyHours: 24,
        defaultVisibility: Visibility.WORKSPACE,
        defaultTags: "",
        shortcutPrefix: "",
      });
      loadRssFeeds();
    } catch (error: any) {
      console.error("Failed to create RSS feed:", error);
      toast.error("Failed to create RSS feed: " + error.message);
    }
  };

  const handleTriggerImport = async (feed: RssFeed) => {
    try {
      setSelectedFeed(feed);
      setShowImportModal(true);
      setImportResult(null);

      const result = await rssImportServiceClient.triggerRssFeedImport({ id: feed.id });
      setImportResult(result);
      toast.success(`Imported ${result.importedCount} shortcuts`);
      loadRssFeeds(); // Refresh the list
    } catch (error: any) {
      console.error("Failed to trigger import:", error);
      toast.error("Failed to import RSS feed: " + error.message);
      setImportResult({ importedCount: 0, errors: [error.message] });
    }
  };

  const handleDeleteFeed = async (feedId: number) => {
    if (!confirm("Are you sure you want to delete this RSS feed?")) {
      return;
    }

    try {
      await rssImportServiceClient.deleteRssFeed({ id: feedId });
      toast.success("RSS feed deleted successfully");
      loadRssFeeds();
    } catch (error: any) {
      console.error("Failed to delete RSS feed:", error);
      toast.error("Failed to delete RSS feed: " + error.message);
    }
  };

  return (
    <Card className="w-full p-6">
      <div className="flex justify-between items-center mb-4">
        <div>
          <Typography level="title-md" className="mb-2">
            RSS Import
          </Typography>
          <Typography level="body-sm" className="text-gray-600 dark:text-gray-400">
            Automatically import links from RSS feeds as shortcuts
          </Typography>
        </div>
        <Button
          variant="outlined"
          color="primary"
          startDecorator={<Icon.Plus />}
          onClick={() => setShowCreateModal(true)}
        >
          Add RSS Feed
        </Button>
      </div>

      {loading ? (
        <div className="flex justify-center items-center py-8">
          <Icon.Loader className="w-6 h-6 animate-spin" />
          <Typography level="body-sm" className="ml-2">Loading RSS feeds...</Typography>
        </div>
      ) : rssFeeds.length === 0 ? (
        <div className="text-center py-8">
          <Icon.Rss className="w-12 h-12 mx-auto mb-4 text-gray-400" />
          <Typography level="body-md" className="text-gray-500 mb-4">
            No RSS feeds configured
          </Typography>
          <Typography level="body-sm" className="text-gray-400">
            Add an RSS feed to automatically import links as shortcuts
          </Typography>
        </div>
      ) : (
        <div className="space-y-4">
          {rssFeeds.map((feed) => (
            <div key={feed.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div className="flex justify-between items-start mb-2">
                <div className="flex-1">
                  <Typography level="title-sm" className="mb-1">
                    {feed.title || feed.url}
                  </Typography>
                  <Typography level="body-xs" className="text-gray-500 mb-2">
                    {feed.url}
                  </Typography>
                  {feed.description && (
                    <Typography level="body-sm" className="text-gray-600 dark:text-gray-400 mb-2">
                      {feed.description}
                    </Typography>
                  )}
                  <div className="flex items-center space-x-4 text-sm text-gray-500">
                    <span>Auto-import: {feed.autoImport ? "On" : "Off"}</span>
                    <span>Frequency: {feed.importFrequencyHours}h</span>
                    <span>Imported: {feed.totalImported}</span>
                    {feed.defaultTags.length > 0 && (
                      <span>Tags: {feed.defaultTags.join(", ")}</span>
                    )}
                  </div>
                </div>
                <div className="flex space-x-2">
                  <Button
                    size="sm"
                    variant="outlined"
                    color="primary"
                    onClick={() => handleTriggerImport(feed)}
                  >
                    Import Now
                  </Button>
                  <Button
                    size="sm"
                    variant="outlined"
                    color="danger"
                    onClick={() => handleDeleteFeed(feed.id)}
                  >
                    <Icon.Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create RSS Feed Modal */}
      <Modal open={showCreateModal} onClose={() => setShowCreateModal(false)}>
        <ModalDialog className="w-full max-w-md">
          <ModalClose />
          <Typography level="h4" className="mb-4">
            Add RSS Feed
          </Typography>

          <div className="space-y-4">
            <Input
              placeholder="Feed Title (optional)"
              value={createForm.title}
              onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
            />

            <Input
              placeholder="RSS Feed URL *"
              value={createForm.url}
              onChange={(e) => setCreateForm({ ...createForm, url: e.target.value })}
              required
            />

            <Textarea
              placeholder="Description (optional)"
              value={createForm.description}
              onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
              minRows={2}
            />

            <Input
              placeholder="Default tags (comma-separated)"
              value={createForm.defaultTags}
              onChange={(e) => setCreateForm({ ...createForm, defaultTags: e.target.value })}
            />

            <Input
              placeholder="Shortcut prefix (optional)"
              value={createForm.shortcutPrefix}
              onChange={(e) => setCreateForm({ ...createForm, shortcutPrefix: e.target.value })}
            />

            <Select
              value={createForm.defaultVisibility}
              onChange={(_, value) => setCreateForm({ ...createForm, defaultVisibility: value as Visibility })}
            >
              <Option value={Visibility.WORKSPACE}>Workspace</Option>
              <Option value={Visibility.PUBLIC}>Public</Option>
            </Select>

            <div className="flex items-center justify-between">
              <Typography level="body-sm">Auto-import</Typography>
              <Switch
                checked={createForm.autoImport}
                onChange={(e) => setCreateForm({ ...createForm, autoImport: e.target.checked })}
              />
            </div>

            <Input
              type="number"
              placeholder="Import frequency (hours)"
              value={createForm.importFrequencyHours}
              onChange={(e) => setCreateForm({ ...createForm, importFrequencyHours: parseInt(e.target.value) || 24 })}
              slotProps={{ input: { min: 1, max: 168 } }}
            />

            <div className="flex justify-end space-x-2">
              <Button variant="outlined" onClick={() => setShowCreateModal(false)}>
                Cancel
              </Button>
              <Button
                variant="solid"
                onClick={handleCreateFeed}
                disabled={!createForm.url}
              >
                Create Feed
              </Button>
            </div>
          </div>
        </ModalDialog>
      </Modal>

      {/* Import Result Modal */}
      <Modal open={showImportModal} onClose={() => setShowImportModal(false)}>
        <ModalDialog className="w-full max-w-md">
          <ModalClose />
          <Typography level="h4" className="mb-4">
            Import Results
          </Typography>

          {importResult ? (
            <div className="space-y-4">
              <Typography level="body-md">
                Successfully imported {importResult.importedCount} shortcuts
              </Typography>

              {importResult.importedShortcuts && importResult.importedShortcuts.length > 0 && (
                <div>
                  <Typography level="body-sm" className="font-semibold mb-2">
                    Created shortcuts:
                  </Typography>
                  <ul className="text-sm text-gray-600 space-y-1">
                    {importResult.importedShortcuts.map((shortcut: string, index: number) => (
                      <li key={index}>• {shortcut}</li>
                    ))}
                  </ul>
                </div>
              )}

              {importResult.errors && importResult.errors.length > 0 && (
                <div>
                  <Typography level="body-sm" className="font-semibold mb-2 text-red-600">
                    Errors:
                  </Typography>
                  <ul className="text-sm text-red-600 space-y-1">
                    {importResult.errors.map((error: string, index: number) => (
                      <li key={index}>• {error}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ) : (
            <div className="flex justify-center items-center py-8">
              <Icon.Loader className="w-6 h-6 animate-spin" />
              <Typography level="body-sm" className="ml-2">Importing...</Typography>
            </div>
          )}
        </ModalDialog>
      </Modal>
    </Card>
  );
};

export default RssImportSection;