import { EnvelopeIcon } from '@heroicons/react/24/solid';
import PropTypes from 'prop-types';
import { useEffect, useState } from 'react';

function Checkout({ item, name, email, date, returned }) {
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
                <p className="text-base mb-2">{returned ? "Not Returned" : "Returned"}</p>
            </div>
        </div>
    )
}

Checkout.propTypes = {
    item: PropTypes.string,
    name: PropTypes.string,
    email: PropTypes.string,
    date: PropTypes.string,
    returned: PropTypes.bool
}

function CheckoutDisplay() {
    const [checkouts, setCheckouts] = useState(Array([]))
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        fetch('http://localhost:8080/checkout')
        .then(response => {
                if(response.ok) {
                    return response.json()
                }
                throw response
            })
        .then(data => {
                console.log(data)
                setCheckouts(data['Checkouts'])
                console.log(checkouts)
            })
        .catch(error => {
                console.log(error)
                setError(error)
            })
        .finally(() => {
                setLoading(false)
            })
    },[])

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

    return(
        <div className="m-2">
            <div className="grid sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {
                    checkouts.map(checkout => {
                        let date = new Date(checkout.Date)
                        return(
                            <Checkout key={checkout.Id} item={checkout.ItemName} name={checkout.Name} email={checkout.Email} date={date.toLocaleDateString()} returned={checkout.Returned}/>
                        )
                    })
                }
            </div>
        </div>
    )
}

export default CheckoutDisplay
