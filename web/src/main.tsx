import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";
import ItemsView from "./views/ItemsView";
import ItemView from "./views/ItemView";
import HomeView from "./views/HomeView";

const root = document.getElementById('root');

createRoot(root!).render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<HomeView />} />
      <Route path="/items" element={<ItemsView />} />
      <Route path="/items/:itemid" element={<ItemView />} />
    </Routes>
  </BrowserRouter>
);

