type CardProps = {
  children: React.ReactNode;
};

const Card: React.FC<CardProps> = ({ children }) => {
  return (
    <div className="shadow-2xl bg-white rounded overflow-hidden py-8 px-5 lg:py-14 md:py-14 sm:py-8 lg:px-12 md:px-12 sm:px-5">
      {children}
    </div>
  );
};

export default Card;
