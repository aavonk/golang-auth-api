import React from "react";
import styles from "./input.module.css";
import classnames from "classnames/bind";

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
  size?: "medium" | "large";
}

const cx = classnames.bind(styles);

const Input = React.forwardRef<HTMLInputElement, TextInputDefaultProps>(
  function Input({ type = "text", size = "medium", ...props }, ref) {
    return (
      <input
        ref={ref}
        className={cx(styles.root, {
          large: size === "large",
          medium: size === "medium",
        })}
        type={type}
        {...props}
      />
    );
  }
);

export default Input;
