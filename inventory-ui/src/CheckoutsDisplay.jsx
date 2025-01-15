import { EnvelopeIcon } from '@heroicons/react/24/solid';
import PropTypes from 'prop-types';
import { useEffect, useState } from 'react';
import Modal from './utilities/Modal'

function ReturnModal ({ name, returnFunc, open, onClose }) {
    return (
        <Modal title={`Return ${name}`} open={open} onClose={onClose} >
            <>
                <p>Are you sure you want to mark as returned?</p>
                <button className="btn btn-danger mt-1 justi" onClick={() => {returnFunc(); onClose()}}>Confirm</button>
            </>
        </Modal>
    )

}

ReturnModal.propTypes = {
    name: PropTypes.string,
    returnFunc: PropTypes.func,
    open: PropTypes.bool,
    onClose: PropTypes.func,
}

function Checkout({ item, name, email, date, returned, returnFunc}) {
    const [open, setOpen] = useState(false)
    return (
        <div className="mt-2 mb-2">
            <div className="bg-ash_gray text-slate-900 text-center p-2 rounded-xl">
                <h1 className="text-xl font-bold mt-2 mb-2">{item}</h1>
                <p className="text-base mb-2 text-center">{name}</p>
                <p className="text-base mb-2 flex justify-center">
                    {email}
                    <a href={`mailto:${email}`} className="text-base">
                        <EnvelopeIcon className='size-6 ps-1'/>
                    </a>
                </p>
                <p className="text-pase mb-2">Checkout Date: {date}</p>
                <p className="text-base mb-2">{returned ? "Returned" : "Not Returned"}</p>
                <button className="btn btn-danger" onClick={() => setOpen(true)}>Mark Returned</button>
                <ReturnModal name={item} returnFunc={returnFunc} open={open} onClose={() => setOpen(false)} />
            </div>
        </div>
    )
}

Checkout.propTypes = {
    item: PropTypes.string,
    name: PropTypes.string,
    email: PropTypes.string,
    date: PropTypes.string,
    returned: PropTypes.bool,
    returnFunc: PropTypes.func
}

function CheckoutDisplay() {
    const [checkouts, setCheckouts] = useState(Array([]))
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    const returnItem = (id) => {
        fetch('http://localhost:8080/checkout?id=' + id, {method: "PUT"})
        .then(response => {
                if(response.ok) {
                    console.log("success")
                    fetchData()
                    return
                }
                alert("Could not mark item as returned")
            })
    }

    function fetchData() {
        fetch('http://localhost:8080/checkout')
        .then(response => {
                if(response.ok) {
                    return response.json()
                }
                throw response
            })
        .then(data => {
                console.log(data)
                setCheckouts(data['checkouts'])
                console.log(checkouts)
            })
        .catch(error => {
                console.log(error)
                setError(error)
            })
        .finally(() => {
                setLoading(false)
            })
    }

    useEffect(() => {
        fetchData()
    }, [])

    if (loading) return(
        <div className="text-center">
            <h1 className="text-4xl font-bold bg-red-700 align-middle pt-10 pb-10">loading...</h1>
        </div>
    )
    if (error) return (
        <div className="text-4xl font-bold bg-red-700 align-middle text-center pt-10 pb-10">
            <h1>Error! please try reloading</h1>
        </div>
    )

    if (!checkouts) {
        return(
            <div className="text-4xl font-bold align-middle text-center pt-10 pb-10">
                <h1>No Checkouts</h1>
            </div>

        )
    }

    return(
        <div className="m-2">
            <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {
                    checkouts.map(checkout => {
                        let date = new Date(checkout.date)
                        return(
                            <Checkout
                                key={checkout.id}
                                item={checkout.itemName}
                                name={checkout.name}
                                email={checkout.email}
                                date={date.toLocaleDateString()}
                                returned={checkout.returned}
                                returnFunc={()=>returnItem(checkout.id)}
                            />
                        )
                    })
                }
            </div>
        </div>
    )
}

export default CheckoutDisplay
