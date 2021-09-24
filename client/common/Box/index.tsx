type BoxProps = {
  margin?: {
    top?: number;
    bottom?: number;
    left?: number;
    right?: number;
  };
  children: React.ReactNode;
};

const Box: React.FC<BoxProps> = ({ margin, children }) => {
  return (
    <div
      style={{
        marginTop: margin?.top,
        marginBottom: margin?.bottom,
        marginLeft: margin?.left,
        marginRight: margin?.right,
      }}
    >
      {children}
    </div>
  );
};

export default Box;
