import React from "react";

export type LabelProps = React.LabelHTMLAttributes<HTMLLabelElement>;

export const Label: React.FC<LabelProps> = ({ className = "", ...props }) => (
  <label className={`block text-sm font-medium text-gray-700 ${className}`} {...props} />
);

export default Label;
