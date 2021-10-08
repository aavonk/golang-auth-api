import React from 'react';
import classnames from 'classnames';

interface Props<C extends React.ElementType> {
  children: React.ReactNode;
  as?: C;
  fullWidth?: boolean;
  type?: 'button' | 'submit';
  onClick?: () => void;
  disabled?: boolean;
}

type ButtonProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithoutRef<C>, keyof Props<C>>;

const Button = <C extends React.ElementType = 'button'>({
  children,
  className,
  fullWidth = false,
  disabled,
  as,
  ...other
}: ButtonProps<C>) => {
  const Component = as || 'button';
  return (
    <Component
      {...other}
      disabled={disabled}
      className={classnames(
        'h-11 inline-flex items-center justify-center rounded text-white py-2 px-4 no-underline',
        {
          'w-full': fullWidth,
          'bg-purple-300 pointer-events-none': disabled,
          'bg-purple-500': !disabled,
        },
        className,
      )}
    >
      {children}
    </Component>
  );
};

export { Button };
