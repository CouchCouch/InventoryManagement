import dragonBeaver from '../assets/dragon-beaver.png';

export const DragonBeavers = () => {
  return(
    <>
      <img src={ dragonBeaver } className="absolute right-4 top-4 h-12" />
      <img src={ dragonBeaver } className="absolute left-4 top-4 h-12" />
    </>
  )
}
