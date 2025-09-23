import { Card, Typography, Button } from "@mui/joy";
import { Link } from "react-router-dom";
import AccessTokenSection from "@/components/setting/AccessTokenSection";
import AccountSection from "@/components/setting/AccountSection";
import PreferenceSection from "@/components/setting/PreferenceSection";
import UserSummarySection from "@/components/setting/UserSummarySection";
import Icon from "@/components/Icon";

const Setting: React.FC = () => {
  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 py-6 flex flex-col justify-start items-start gap-y-12">
      <AccountSection />
      <UserSummarySection />
      <AccessTokenSection />
      <PreferenceSection />

      <Card className="w-full p-6">
        <Typography level="title-md" className="mb-3">
          Data Management
        </Typography>
        <Typography level="body-sm" className="mb-4 text-gray-600 dark:text-gray-400">
          Import and manage your data
        </Typography>
        <Link to="/admin/import">
          <Button
            variant="outlined"
            color="primary"
            startDecorator={<Icon.Upload />}
          >
            Import Bookmarks
          </Button>
        </Link>
      </Card>
    </div>
  );
};

export default Setting;
