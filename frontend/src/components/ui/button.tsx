import React from "react";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "outline" | "destructive";
}

const variantStyles: Record<ButtonProps["variant"], string> = {
  primary:
    "bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500",
  outline:
    "border border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-indigo-500",
  destructive:
    "bg-red-600 text-white hover:bg-red-700 focus:ring-red-500",
};

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className = "", variant = "primary", ...props }, ref) => (
    <button
      ref={ref}
      className={`inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium shadow-sm focus:outline-none focus:ring-2 ${variantStyles[variant]} ${className}`}
      {...props}
    />
  )
);

Button.displayName = "Button";

export default Button;
