import { requireAuthentication } from "hoc/withAuth";
import { GetServerSideProps, NextPage } from "next";

const Home: NextPage = () => {
  return <div>You are logged in!</div>;
};

export default Home;

//@ts-ignore
export const getServerSideProps: GetServerSideProps = requireAuthentication(
  async (_ctx) => {
    return {
      props: {},
    };
  }
);
