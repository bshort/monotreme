import { useState } from "react";
import Icon from "./Icon";
import LinkFavicon from "./LinkFavicon";

interface Props {
  customIcon?: string;
  url?: string;
  className?: string;
}

const CustomIcon = (props: Props) => {
  const { customIcon, url, className = "w-full h-auto rounded" } = props;
  const [hasError, setHasError] = useState<boolean>(false);

  const handleImgError = () => {
    setHasError(true);
  };

  if (customIcon && !hasError) {
    return (
      <img
        className={className}
        src={customIcon}
        decoding="async"
        loading="lazy"
        onError={handleImgError}
        alt="Custom icon"
      />
    );
  }

  if (url) {
    return <LinkFavicon url={url} />;
  }

  return <Icon.Earth className={className.replace("rounded", "")} strokeWidth={1.5} />;
};

export default CustomIcon;