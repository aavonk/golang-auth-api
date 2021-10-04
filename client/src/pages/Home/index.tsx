import { Link } from "react-router-dom";
import Box from "../../common/Box";
import { Button } from "../../common/Button";

function Home() {
  return (
    <div
      style={{
        width: "100%",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Box margin={{ right: 10 }}>
        <Button as={Link} to="/login">
          Login
        </Button>
      </Box>
      <Button as={Link} to="/register">
        Sign up
      </Button>
    </div>
  );
}

export default Home;
