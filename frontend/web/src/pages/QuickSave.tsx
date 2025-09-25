import { Alert, Button, Card, Typography, Input, CircularProgress } from "@mui/joy";
import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import Icon from "@/components/Icon";
import useShortcutStore from "@/stores/shortcut";
import { useUserStore } from "@/stores";
import { Visibility } from "@/types/proto/api/v1/common";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";

const QuickSave = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");
  const [name, setName] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState(false);
  const [savedShortcut, setSavedShortcut] = useState<Shortcut | null>(null);
  const [autoCloseTimer, setAutoCloseTimer] = useState<number>(0);

  const shortcutStore = useShortcutStore();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();

  useEffect(() => {
    // Get URL and title from query parameters
    const urlParam = searchParams.get("url");
    const titleParam = searchParams.get("title");

    if (urlParam) {
      setUrl(decodeURIComponent(urlParam));
    }

    if (titleParam) {
      const decodedTitle = decodeURIComponent(titleParam);
      setTitle(decodedTitle);

      // Generate a shortcut name from the title
      const generatedName = decodedTitle
        .toLowerCase()
        .replace(/\s+/g, '-')
        .replace(/[^\w-]/g, '')
        .substring(0, 32);
      setName(generatedName);
    }
  }, [searchParams]);

  // Redirect to sign-in if not authenticated
  useEffect(() => {
    if (!currentUser.id) {
      navigate("/auth/signin");
    }
  }, [currentUser, navigate]);

  const handleSave = async () => {
    if (!title.trim() || !url.trim() || !name.trim()) {
      setError("Please fill in all fields");
      return;
    }

    setIsSaving(true);
    setError("");

    try {
      const shortcutCreate: Shortcut = {
        id: 0,
        uuid: "",
        creatorId: currentUser.id,
        createdTime: undefined,
        updatedTime: undefined,
        name: name.trim(),
        link: url.trim(),
        title: title.trim(),
        tags: [],
        description: "",
        visibility: currentUser.defaultVisibility || Visibility.WORKSPACE,
        viewCount: 0,
        ogMetadata: undefined,
      };

      const createdShortcut = await shortcutStore.createShortcut(shortcutCreate);
      setSavedShortcut(createdShortcut);
      setSuccess(true);

      // Start auto-close countdown
      let countdown = 5;
      setAutoCloseTimer(countdown);
      const interval = setInterval(() => {
        countdown -= 1;
        setAutoCloseTimer(countdown);
        if (countdown <= 0) {
          clearInterval(interval);
          window.close();
        }
      }, 1000);
    } catch (err: any) {
      setError(err.message || "Failed to save bookmark");
    } finally {
      setIsSaving(false);
    }
  };

  const handleGoHome = () => {
    navigate("/");
  };

  const handleViewShortcut = () => {
    if (savedShortcut) {
      navigate(`/s/${savedShortcut.name}`);
    }
  };

  if (!currentUser.id) {
    return null; // Will redirect to sign-in
  }

  return (
    <div className="w-full px-4 py-4 flex flex-col gap-y-4 min-h-screen">
      <div className="flex justify-center">
        <div className="flex items-center gap-3">
          <img src="../assets/monotreme.png" alt="Monotreme Logo" className="w-16 h-16" />
          <div>
            <Typography level="h3" className="mb-0">
              Quick Save
            </Typography>
            <Typography level="body-sm" className="text-gray-600 dark:text-gray-400">
              Save to Mon.otre.me
            </Typography>
          </div>
        </div>
      </div>

      {!success ? (
        <Card className="w-full p-4">
          <div className="space-y-3">
            <div>
              <Typography level="title-sm" className="mb-1">
                Title
              </Typography>
              <Input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Page title"
                disabled={isSaving}
                fullWidth
                size="sm"
              />
            </div>

            <div>
              <Typography level="title-sm" className="mb-1">
                URL
              </Typography>
              <Input
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="Page URL"
                disabled={isSaving}
                fullWidth
                size="sm"
              />
            </div>

            <div>
              <Typography level="title-sm" className="mb-1">
                Shortcut Name
              </Typography>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Short name"
                disabled={isSaving}
                fullWidth
                size="sm"
              />
              <Typography level="body-xs" className="mt-1 text-gray-600">
                URL: /s/{name}
              </Typography>
            </div>

            <div className="flex gap-2 pt-2">
              <Button
                variant="solid"
                color="primary"
                onClick={handleSave}
                disabled={isSaving}
                startDecorator={isSaving ? <CircularProgress size="sm" /> : <Icon.BookmarkPlus />}
                fullWidth
                size="sm"
              >
                {isSaving ? "Saving..." : "Save"}
              </Button>

              <Button
                variant="outlined"
                color="neutral"
                onClick={() => window.close()}
                disabled={isSaving}
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </div>
        </Card>
      ) : (
        <Card className="w-full p-4">
          <Alert variant="solid" color="success" startDecorator={<Icon.CheckCircle />} className="mb-3">
            Saved successfully! {autoCloseTimer > 0 && `(Closing in ${autoCloseTimer}s)`}
          </Alert>

          <div className="space-y-3 text-center">
            <div>
              <Typography level="title-sm" className="mb-1">
                {savedShortcut?.title}
              </Typography>
              <Typography level="body-xs" className="text-gray-600 break-all">
                {savedShortcut?.link}
              </Typography>
            </div>

            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-2">
              <Typography level="body-xs" className="font-medium mb-1">
                Shortcut URL:
              </Typography>
              <Typography level="body-xs" className="font-mono text-blue-700 dark:text-blue-300 break-all">
                /s/{savedShortcut?.name}
              </Typography>
            </div>

            <div className="flex gap-2 justify-center">
              <Button
                variant="solid"
                color="primary"
                onClick={handleViewShortcut}
                startDecorator={<Icon.ExternalLink />}
                size="sm"
              >
                View
              </Button>

              <Button
                variant="outlined"
                color="neutral"
                onClick={() => window.close()}
                size="sm"
              >
                Close
              </Button>
            </div>
          </div>
        </Card>
      )}

      {error && (
        <Alert variant="solid" color="danger" startDecorator={<Icon.AlertTriangle />}>
          {error}
        </Alert>
      )}
    </div>
  );
};

export default QuickSave;