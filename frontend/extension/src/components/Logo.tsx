import classNames from "classnames";

interface Props {
  className?: string;
}

const Logo = ({ className }: Props) => {
  return (
    <div className={classNames("w-8 h-8 rounded-lg overflow-hidden", className)}>
      <img src="../assets/monotreme.png" alt="Monotreme Logo" className="max-w-full max-h-full" />
    </div>
  );
};

export default Logo;
