import { Link } from "react-router";
import { Header } from "../components/Header";

export default function HomeView() {
  return (
    <>
      <Header />
      <div className="flex items-center justify-center">
        <Link to="/items" >
          <button className="btn btn-create">View Items</button>
        </Link>
      </div>
    </>
  )
}
