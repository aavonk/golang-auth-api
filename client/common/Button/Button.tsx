import React from "react";
import styled from "styled-components";
import Link from "next/link";

const StyledButton = styled.button`
  min-height: 44px;
  display: inline-flex;
  border-radius: 4px;
  background-color: #654ef5;
  color: #fff;
  align-items: center;
  justify-content: center;
  padding: 8px 16px;
  vertical-align: middle;
  text-decoration: none;
  cursor: pointer;
  & > a {
    min-height: 44px;
    display: inline-flex;
    border-radius: 4px;
    background-color: #654ef5;
    color: #fff;
    align-items: center;
    justify-content: center;
    padding: 8px 16px;
    vertical-align: middle;
    text-decoration: none;
    cursor: pointer;
  }
`;

interface Props<C extends React.ElementType> {
  children: React.ReactNode;
  as?: C;
}

type ButtonProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithoutRef<C>, keyof Props<C>>;

const Button = <C extends React.ElementType = "button">({
  children,
  className,
  as,
  ...other
}: ButtonProps<C>) => {
  const Component = as || "button";
  return (
    //@ts-ignore
    <StyledButton as={Component} {...other}>
      {children}
    </StyledButton>
  );
};

export { Button };
