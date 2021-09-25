type Quadrant = {
  top?: number;
  bottom?: number;
  left?: number;
  right?: number;
};

type BoxProps = {
  margin?: Quadrant;
  padding?: Quadrant;
  children: React.ReactNode;
};

const Box: React.FC<BoxProps> = ({ margin, children, padding }) => {
  return (
    <div
      style={{
        marginTop: margin?.top,
        marginBottom: margin?.bottom,
        marginLeft: margin?.left,
        marginRight: margin?.right,
        paddingTop: padding?.top,
        paddingBottom: padding?.bottom,
        paddingRight: padding?.right,
        paddingLeft: padding?.left,
      }}
    >
      {children}
    </div>
  );
};

export default Box;
