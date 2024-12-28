import { XMarkIcon } from "@heroicons/react/24/solid";

export default function Modal({ open, onClose, title, children }) {
    return (
        <div
            onClick={onClose}
            className={
                `fixed inset-0 flex justify-center items-center
                 transition-colors ${open ? "visible bg-black/20" : "invisible"}
            `}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                className={`
                    bg-lavender rounded-xl shadow p-6 transition-all
                    ${open ? "scale-100 opacity-100" : "scale-125 opacity-0"}
                `}>
                    <div className="block text-center">
                        <h3 className="text-lg font-black text-gray-800 mb-1">{title}</h3>
                        <button
                            onClick={onClose}
                            className="btn absolute top-2 right-2 p-1 rounded-lg text-gray-400 bg-lavender hover:bg-red-600 hover:text-white"
                        >
                            <XMarkIcon className="size-6"/>
                        </button>
                    </div>
                    {children}
                </div>
        </div>
    )
}
