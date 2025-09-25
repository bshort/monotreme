import { Button } from "@mui/joy";
import { useRef, useState } from "react";
import { toast } from "react-hot-toast";
import Icon from "./Icon";
import CustomIcon from "./CustomIcon";

interface Props {
  value?: string;
  onChange: (iconData: string) => void;
  fallbackUrl?: string;
  className?: string;
}

const IconUpload = (props: Props) => {
  const { value, onChange, fallbackUrl, className = "w-16 h-16" } = props;
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState<boolean>(false);

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      toast.error("Please select a valid image file (PNG, JPG, ICO, etc.)");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("File size must be less than 5MB");
      return;
    }

    setIsUploading(true);

    try {
      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result as string;
        onChange(result);
        setIsUploading(false);
        toast.success("Icon uploaded successfully");
      };
      reader.onerror = () => {
        setIsUploading(false);
        toast.error("Failed to read file");
      };
      reader.readAsDataURL(file);
    } catch (error) {
      setIsUploading(false);
      toast.error("Failed to upload icon");
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleRemoveIcon = () => {
    onChange("");
    toast.success("Icon removed");
  };

  return (
    <div className="flex flex-col items-center space-y-2">
      <div className={`${className} border border-dashed border-gray-300 dark:border-gray-600 rounded-lg flex items-center justify-center overflow-hidden`}>
        {value || fallbackUrl ? (
          <CustomIcon
            customIcon={value}
            url={fallbackUrl}
            className="w-full h-full object-cover"
          />
        ) : (
          <Icon.ImageIcon className="w-8 h-8 text-gray-400" />
        )}
      </div>

      <div className="flex space-x-2">
        <Button
          size="sm"
          variant="outlined"
          onClick={handleUploadClick}
          loading={isUploading}
          disabled={isUploading}
        >
          <Icon.Upload className="w-4 h-4 mr-1" />
          {value ? "Change" : "Upload"}
        </Button>

        {value && (
          <Button
            size="sm"
            variant="outlined"
            color="danger"
            onClick={handleRemoveIcon}
            disabled={isUploading}
          >
            <Icon.Trash2 className="w-4 h-4 mr-1" />
            Remove
          </Button>
        )}
      </div>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        className="hidden"
      />
    </div>
  );
};

export default IconUpload;