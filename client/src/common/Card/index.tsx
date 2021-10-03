import classnames from "classnames/bind";
import styles from "./card.module.css";

const cx = classnames.bind(styles);

type CardProps = {
  children: React.ReactNode;
};

const Card: React.FC<CardProps> = ({ children }) => {
  return <div className={cx(styles.card_root)}>{children}</div>;
};

export default Card;
