import React from "react";
import styles from "./button.module.css";
import classnames from "classnames/bind";

interface Props<C extends React.ElementType> {
  children: React.ReactNode;
  as?: C;
  fullWidth?: boolean;
  type?: "button" | "submit";
  onClick?: () => void;
}

type ButtonProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithoutRef<C>, keyof Props<C>>;

const cx = classnames.bind(styles);

const Button = <C extends React.ElementType = "button">({
  children,
  className,
  fullWidth = false,
  as,
  ...other
}: ButtonProps<C>) => {
  const Component = as || "button";
  return (
    <Component
      {...other}
      className={cx(styles.button_root, {
        fullWidth: fullWidth,
      })}
    >
      {children}
    </Component>
  );
};

export { Button };
