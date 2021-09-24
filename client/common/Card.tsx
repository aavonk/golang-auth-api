import styled from "styled-components";

const StyledCard = styled.div`
  border-radius: 4px;
  box-shadow: 0 15px 35px 0 rgb(60 66 87 / 8%), 0 5px 15px 0 rgb(0 0 0 / 12%);
  overflow: hidden;
  background-color: #fff;
  padding: 56px 48px;

  @media (max-width: 880px) {
    padding: 32px 20px;
  }
`;

type CardProps = {
  children: React.ReactNode;
};
const Card: React.FC<CardProps> = ({ children }) => {
  return <StyledCard>{children}</StyledCard>;
};

export default Card;
