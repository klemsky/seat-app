import React, { useState, useEffect } from 'react';

function App() {
  const [seatMapData, setSeatMapData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedSeat, setSelectedSeat] = useState(null);

  useEffect(() => {
    const fetchSeatMap = async () => {
      try {
        const response = await fetch('/api/seatmap');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setSeatMapData(data);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };

    fetchSeatMap();
  }, []);

  const handleSeatClick = (seat) => {
    if (seat.available && !seat.freeOfCharge && seat.storefrontSlotCode === "SEAT") {
      setSelectedSeat(seat);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100 font-inter">
        <div className="text-xl font-semibold text-gray-700">Loading seat map...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100 font-inter">
        <div className="text-xl font-semibold text-red-600">Error: {error}</div>
      </div>
    );
  }

  if (!seatMapData || seatMapData.seatsItineraryParts.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100 font-inter">
        <div className="text-xl font-semibold text-gray-700">No seat map data available.</div>
      </div>
    );
  }

  const { segmentSeatMaps } = seatMapData.seatsItineraryParts[0];
  const { segment, passengerSeatMaps } = segmentSeatMaps[0];
  const { passenger, seatMap } = passengerSeatMaps[0];

  return (
    <div className="min-h-screen bg-gray-100 p-4 font-inter">
      <div className="max-w-6xl mx-auto bg-white shadow-lg rounded-xl p-8 space-y-8">
        {/* Header */}
        <h1 className="text-4xl font-extrabold text-center text-indigo-800 mb-8">
          Flight Seat Map
        </h1>

        {/* Flight Details Section */}
        <div className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white p-6 rounded-lg shadow-md flex justify-between items-center">
          <div>
            <p className="text-lg font-bold">
              Flight {segment.flight.airlineCode}{segment.flight.flightNumber}
            </p>
            <p className="text-sm">
              {segment.origin} to {segment.destination}
            </p>
          </div>
          <div className="text-right">
            <p className="text-lg font-bold">
              Departure: {new Date(segment.departure).toLocaleString()}
            </p>
            <p className="text-sm">
              Arrival: {new Date(segment.arrival).toLocaleString()}
            </p>
          </div>
        </div>

        {/* Passenger Information */}
        <div className="bg-blue-50 border border-blue-200 p-6 rounded-lg shadow-md">
          <h2 className="text-2xl font-bold text-blue-800 mb-4">Passenger Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-gray-700">
            <p>
              <span className="font-semibold">Name:</span> {passenger.passengerDetails.firstName}{" "}
              {passenger.passengerDetails.lastName}
            </p>
            <p>
              <span className="font-semibold">Type:</span> {passenger.passengerInfo.type}
            </p>
            <p>
              <span className="font-semibold">Date of Birth:</span>{" "}
              {passenger.passengerInfo.dateOfBirth}
            </p>
            <p>
              <span className="font-semibold">Gender:</span> {passenger.passengerInfo.gender}
            </p>
            <p>
              <span className="font-semibold">Frequent Flyer:</span>{" "}
              {passenger.preferences.frequentFlyer[0]?.airline}{" "}
              {passenger.preferences.frequentFlyer[0]?.number}
            </p>
          </div>
        </div>

        {/* Seat Map */}
        <div className="overflow-x-auto">
          <div className="bg-gray-50 border border-gray-200 p-6 rounded-lg shadow-md min-w-max">
            <h2 className="text-2xl font-bold text-gray-800 mb-4 text-center">
              Select Your Seat ({seatMap.aircraft})
            </h2>
            <div className="flex justify-center mb-4 space-x-4">
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 bg-green-400 rounded-md"></div>
                <span>Available</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 bg-red-400 rounded-md"></div>
                <span>Occupied/Unavailable</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 bg-yellow-400 rounded-md"></div>
                <span>Premium/Extra Cost</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 bg-blue-400 rounded-md border-2 border-blue-600"></div>
                <span>Selected</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 bg-gray-300 rounded-md"></div>
                <span>Aisle/Blank/Wing</span>
              </div>
            </div>

            {seatMap.cabins.map((cabin, cabinIndex) => (
              <div key={cabinIndex} className="mb-8">
                <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
                  {cabin.deck} Cabin
                </h3>
                {/* Column Headers */}
                <div className="flex justify-center text-gray-600 font-bold mb-2">
                  {cabin.seatColumns.map((col, colIdx) => (
                    <div
                      key={colIdx}
                      className={`w-10 h-10 flex items-center justify-center rounded-md ${col === "AISLE" || col.includes("SIDE") || col === "BLANK" || col === "WING"
                        ? "bg-transparent text-gray-400"
                        : "text-gray-700"
                        }`}
                    >
                      {col === "LEFT_SIDE" ? "<" : col === "RIGHT_SIDE" ? ">" : col === "AISLE" ? " " : col === "BLANK" ? " " : col === "WING" ? "~" : col}
                    </div>
                  ))}
                </div>

                {/* Seat Rows */}
                {cabin.seatRows.map((row, rowIndex) => (
                  <div key={rowIndex} className="flex justify-center mb-2 items-center">
                    {/* Row Number */}
                    <div className="w-10 h-10 flex items-center justify-center font-bold text-gray-700 mr-2">
                      {row.rowNumber}
                    </div>
                    {/* Seats in the row */}
                    {row.seats.map((seat, seatIndex) => {
                      let seatClass = "w-10 h-10 flex items-center justify-center rounded-md cursor-pointer transition-all duration-200";
                      let seatContent = seat.code || "";
                      let tooltipContent = "";

                      if (seat.storefrontSlotCode === "AISLE") {
                        seatClass = "w-10 h-10 flex items-center justify-center text-gray-400";
                        seatContent = " ";
                      } else if (seat.storefrontSlotCode === "BLANK" || seat.storefrontSlotCode === "WING") {
                        seatClass = "w-10 h-10 flex items-center justify-center text-gray-400";
                        seatContent = seat.storefrontSlotCode === "WING" ? "~" : " ";
                      } else if (seat.available && !seat.freeOfCharge) {
                        seatClass += " bg-yellow-400 hover:bg-yellow-500 active:scale-95";
                        tooltipContent = `Price: ${seat.total?.alternatives[0][0].amount} ${seat.total?.alternatives[0][0].currency}`;
                      } else if (seat.available && seat.freeOfCharge) {
                        seatClass += " bg-green-400 hover:bg-green-500 active:scale-95";
                        tooltipContent = "Free of charge";
                      } else {
                        seatClass += " bg-red-400 cursor-not-allowed opacity-70";
                        tooltipContent = "Unavailable";
                      }

                      if (selectedSeat && selectedSeat.code === seat.code) {
                        seatClass += " border-2 border-blue-600 scale-105";
                      }

                      return (
                        <div
                          key={seatIndex}
                          className={seatClass + " mx-0.5"} // Add horizontal margin to separate seats
                          onClick={() => handleSeatClick(seat)}
                          title={tooltipContent}
                        >
                          {seatContent}
                        </div>
                      );
                    })}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>

        {/* Selected Seat Information */}
        {selectedSeat && (
          <div className="bg-purple-50 border border-purple-200 p-6 rounded-lg shadow-md">
            <h2 className="text-2xl font-bold text-purple-800 mb-4">Selected Seat</h2>
            <div className="text-gray-700">
              <p>
                <span className="font-semibold">Seat Number:</span> {selectedSeat.code}
              </p>
              <p>
                <span className="font-semibold">Price:</span>{" "}
                {selectedSeat.total?.alternatives[0][0].amount}{" "}
                {selectedSeat.total?.alternatives[0][0].currency}
              </p>
              <p>
                <span className="font-semibold">Availability:</span>{" "}
                {selectedSeat.available ? "Available" : "Not Available"}
              </p>
              <p>
                <span className="font-semibold">Characteristics:</span>{" "}
                {selectedSeat.seatCharacteristics.join(", ")}
              </p>
              <p>
                <span className="font-semibold">Limitations:</span>{" "}
                {selectedSeat.limitations.length > 0
                  ? selectedSeat.limitations.join(", ")
                  : "None"}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
