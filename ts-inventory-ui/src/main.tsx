import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";
import ItemView from "./views/ItemView";

const root = document.getElementById('root');

createRoot(root!).render(
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<ItemView />} />
            <Route path="/items" element={<ItemView />} />
        </Routes>
    </BrowserRouter>
);

