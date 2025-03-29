
import { useSearchParams } from "react-router";
import { Header } from "../components/Header";
import Item from "../components/SingleItem";
import { CheckoutHistory } from "../components/CheckoutHistory";

export default function ItemView() {
    const [params] = useSearchParams()
    const id = Number(params.get('id'))

    return (
        <>
            <Header />
            <div className="p-2">
                <Item id={id} />
                <CheckoutHistory id={id} />
            </div>
        </>
    )
}
