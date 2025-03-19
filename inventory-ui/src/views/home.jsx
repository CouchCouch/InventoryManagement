import ItemDisplay from './ItemsDisplay'
import Headers from './Headers'
import DragonBeavers from './DragonBeavers'
import CheckoutDisplay from './CheckoutsDisplay'

function Home() {
  return (
    <div>
        <Headers />
        <ItemDisplay />
        <CheckoutDisplay />
        <DragonBeavers />
    </div>
  )
}

export default Home
