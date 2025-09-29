import { Card, Typography, Button } from "@mui/joy";
import toast from "react-hot-toast";
import { Link } from "react-router-dom";
import Icon from "@/components/Icon";
import AccessTokenSection from "@/components/setting/AccessTokenSection";
import AccountSection from "@/components/setting/AccountSection";
import PreferenceSection from "@/components/setting/PreferenceSection";
import RecentActivitySection from "@/components/setting/RecentActivitySection";
import UserSummarySection from "@/components/setting/UserSummarySection";
import { userServiceClient } from "@/grpcweb";
import { useUserStore } from "@/stores";

const Setting: React.FC = () => {
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();

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
        <div className="flex flex-row gap-3">
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
      </Card>
    </div>
  );
};

export default Setting;
