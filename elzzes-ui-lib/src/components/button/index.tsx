import type { ButtonProps } from "./button.types";

const variantClasses = {
  primary:
    "bg-ui-primary text-white hover:bg-ui-primary-hover",
  secondary:
    "bg-gray-100 text-gray-900 hover:bg-gray-200",
  outline:
    "border border-gray-300 bg-transparent text-gray-900 hover:bg-gray-100",
  ghost:
    "bg-transparent text-gray-900 hover:bg-gray-100",
  danger:
    "bg-ui-danger text-white hover:opacity-90",
};

const sizeClasses = {
  sm: "px-3 py-1.5 text-sm",
  md: "px-4 py-2 text-sm",
  lg: "px-5 py-3 text-base",
};

export function Button({
  variant = "primary",
  size = "md",
  className = "",
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={[
        "inline-flex items-center justify-center",
        "rounded-lg font-medium transition-colors",
        "focus-visible:outline-2 focus-visible:outline-offset-2",
        "disabled:pointer-events-none disabled:opacity-50",
        variantClasses[variant],
        sizeClasses[size],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...props}
    >
      {children}
    </button>
  );
}