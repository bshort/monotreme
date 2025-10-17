import { Card, Typography, Button, Modal, ModalDialog, Input } from "@mui/joy";
import { useState } from "react";
import toast from "react-hot-toast";
import { Link } from "react-router-dom";
import Icon from "@/components/Icon";
import AccessTokenSection from "@/components/setting/AccessTokenSection";
import AccountSection from "@/components/setting/AccountSection";
import PreferenceSection from "@/components/setting/PreferenceSection";
import RecentActivitySection from "@/components/setting/RecentActivitySection";
import UserSummarySection from "@/components/setting/UserSummarySection";
import { shortcutServiceClient, userServiceClient } from "@/grpcweb";
import { useUserStore, useShortcutStore } from "@/stores";

const Setting: React.FC = () => {
  const userStore = useUserStore();
  const shortcutStore = useShortcutStore();
  const currentUser = userStore.getCurrentUser();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteConfirmText, setDeleteConfirmText] = useState("");

  const handleExportShortcuts = async () => {
    try {
      // Generate a new access token for export
      const { accessToken } = await userServiceClient.createUserAccessToken({
        id: currentUser.id,
        description: "Export Shortcuts as HTML",
        expiresAt: undefined, // Never expires
      });

      // Create the export URL
      const exportUrl = `${window.location.origin}/export/shortcuts.html?token=${accessToken}`;

      // Trigger download
      window.open(exportUrl, "_blank");
      toast.success("Export started! Your browser will download the bookmark file.");
    } catch (error: any) {
      console.error("Failed to export shortcuts:", error);
      toast.error("Failed to export shortcuts. Please try again.");
    }
  };

  const handleExportShortcutsCSV = async () => {
    try {
      // Generate a new access token for export
      const { accessToken } = await userServiceClient.createUserAccessToken({
        id: currentUser.id,
        description: "Export Shortcuts as CSV",
        expiresAt: undefined, // Never expires
      });

      // Create the export URL
      const exportUrl = `${window.location.origin}/export/shortcuts.csv?token=${accessToken}`;

      // Trigger download
      window.open(exportUrl, "_blank");
      toast.success("Export started! Your browser will download the CSV file.");
    } catch (error: any) {
      console.error("Failed to export shortcuts:", error);
      toast.error("Failed to export shortcuts. Please try again.");
    }
  };

  const handleDeleteAllShortcuts = async () => {
    // Check if confirmation text matches (case-insensitive)
    if (deleteConfirmText.toLowerCase() !== "delete all my shortcuts") {
      toast.error("Please type the confirmation text correctly.");
      return;
    }

    try {
      await shortcutServiceClient.deleteAllShortcuts({});
      await shortcutStore.fetchShortcutList();
      toast.success("All shortcuts have been deleted.");
      setShowDeleteModal(false);
      setDeleteConfirmText("");
    } catch (error: any) {
      console.error("Failed to delete all shortcuts:", error);
      toast.error("Failed to delete all shortcuts. Please try again.");
    }
  };

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 py-6 flex flex-col justify-start items-start gap-y-12">
      <AccountSection />
      <UserSummarySection />
      <RecentActivitySection />
      <AccessTokenSection />
      <PreferenceSection />

      <Card className="w-full p-6">
        <Typography level="title-md" className="mb-3">
          Data Management
        </Typography>
        <Typography level="body-sm" className="mb-4 text-gray-600 dark:text-gray-400">
          Import and export your data
        </Typography>
        <div className="flex flex-row gap-3 mb-6">
          <Link to="/admin/import">
            <Button variant="outlined" color="primary" startDecorator={<Icon.Upload />}>
              Import Bookmarks
            </Button>
          </Link>
          <Button variant="outlined" color="neutral" startDecorator={<Icon.Download />} onClick={handleExportShortcuts}>
            Export Shortcuts as HTML
          </Button>
          <Button variant="outlined" color="neutral" startDecorator={<Icon.Download />} onClick={handleExportShortcutsCSV}>
            Export Shortcuts as CSV
          </Button>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
          <Typography level="title-sm" className="mb-2 text-red-600 dark:text-red-400">
            Danger Zone
          </Typography>
          <Typography level="body-sm" className="mb-3 text-gray-600 dark:text-gray-400">
            Permanently delete all your shortcuts. This action cannot be undone.
          </Typography>
          <Button
            variant="outlined"
            color="danger"
            startDecorator={<Icon.Trash2 />}
            onClick={() => setShowDeleteModal(true)}
          >
            Delete All Shortcuts
          </Button>
        </div>
      </Card>

      {/* Delete Confirmation Modal */}
      <Modal open={showDeleteModal} onClose={() => {
        setShowDeleteModal(false);
        setDeleteConfirmText("");
      }}>
        <ModalDialog>
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-2">
              <Icon.AlertTriangle className="w-6 h-6 text-red-600" />
              <Typography level="h4">Delete All Shortcuts</Typography>
            </div>

            <Typography level="body-md" className="text-gray-700 dark:text-gray-300">
              This will permanently delete <strong>all of your shortcuts</strong>. This action cannot be undone.
            </Typography>

            <Typography level="body-sm" className="text-gray-600 dark:text-gray-400">
              To confirm, please type <strong>delete all my shortcuts</strong> below:
            </Typography>

            <Input
              placeholder="delete all my shortcuts"
              value={deleteConfirmText}
              onChange={(e) => setDeleteConfirmText(e.target.value)}
              autoFocus
            />

            <div className="flex justify-end gap-2 mt-2">
              <Button
                variant="plain"
                color="neutral"
                onClick={() => {
                  setShowDeleteModal(false);
                  setDeleteConfirmText("");
                }}
              >
                Cancel
              </Button>
              <Button
                variant="solid"
                color="danger"
                disabled={deleteConfirmText.toLowerCase() !== "delete all my shortcuts"}
                onClick={handleDeleteAllShortcuts}
              >
                Delete All Shortcuts
              </Button>
            </div>
          </div>
        </ModalDialog>
      </Modal>
    </div>
  );
};

export default Setting;
