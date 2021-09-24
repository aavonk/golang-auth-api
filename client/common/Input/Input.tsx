import React from "react";
import styles from "./input.module.css";

type NativeInputProps = React.AllHTMLAttributes<HTMLInputElement>;

export interface TextInputDefaultProps {
  value: NonNullable<string>;
  onChange: NonNullable<NativeInputProps["onChange"]>;
  onBlur?: NativeInputProps["onBlur"];
  type?: string;
  id?: string;
  autoFocus?: NativeInputProps["autoFocus"];
  autoComplete?: NativeInputProps["autoComplete"];
  "aria-describedby"?: NativeInputProps["aria-describedby"];
  disabled?: boolean;
  "data-testid"?: string;
  placeholder?: string;
  icon?: React.ComponentType<{ className?: string }>;
  name?: string;
}

const Input = React.forwardRef<HTMLInputElement, TextInputDefaultProps>(
  function Input({ type = "text", ...props }, ref) {
    return <input className={styles.root} type={type} {...props} />;
  }
);

export default Input;
