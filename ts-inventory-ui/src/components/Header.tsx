import { Link } from "react-router"

export const Header = () => {
    return(
        <div className="text-center">
            <h1 className="text-5xl font-bold pt-4 pb-4 ps-1" ><Link to={"/"}>RITOC Inventory Management System</Link></h1>
        </div>
    )
}
