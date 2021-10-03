import { Link } from "react-router-dom";
import { Button } from "../../common/Button";

function Home() {
  return (
    <div style={{ width: "100%", height: 400 }}>
      <Button as={Link} to="/login">
        Login
      </Button>
      <Button as={Link} to="/register">
        Sign up
      </Button>
    </div>
  );
}

export default Home;
