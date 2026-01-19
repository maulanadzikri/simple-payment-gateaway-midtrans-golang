import React, { useState, useEffect } from 'react';
import { paymentAPI } from '../services/api';

const History = () => {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchHistory = async () => {
    try {
      const response = await paymentAPI.getHistory();
      console.log('History response:', response.data);
      setTransactions(Array.isArray(response.data) ? response.data : []);
    } catch (err) {
      console.error('History error:', err);
      setError(err.response?.data?.error || err.response?.data?.message || 'Failed to fetch history');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHistory();
  }, []);

  const handlePayNow = (url) => {
    if (url) {
      window.open(url, '_blank');
    } else {
      alert("Payment URL tidak ditemukan untuk transaksi ini.");
    }
  };

  const handleCheckStatus = async (orderId) => {
    try {
      setLoading(true);
      await paymentAPI.getStatus(orderId);
      // Refresh data setelah cek status
      await fetchHistory();
      alert(`Status update requested for ${orderId}`);
    } catch (err) {
      console.error('Check status error:', err);
      alert('Failed to check status');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'settlement':
      case 'success':
        return '#28a745';
      case 'pending':
        return '#ffc107';
      case 'cancel':
      case 'expire':
      case 'failed':
        return '#dc3545';
      default:
        return '#6c757d';
    }
  };

  if (loading && transactions.length === 0) {
    return <div style={{ textAlign: 'center', padding: '2rem' }}>Loading...</div>;
  }

  return (
    <div style={{ maxWidth: '1100px', margin: '0 auto' }}>
      <h1 style={{ marginBottom: '2rem' }}>Transaction History</h1>
      
      {error && (
        <div style={{ padding: '0.75rem', backgroundColor: '#f8d7da', color: '#721c24', borderRadius: '4px', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {transactions.length === 0 ? (
        <div style={{ 
          backgroundColor: 'white', 
          padding: '2rem', 
          borderRadius: '8px', 
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <p>No transactions found.</p>
        </div>
      ) : (
        <div style={{ backgroundColor: 'white', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0,0,0,0.1)', overflow: 'hidden' }}>
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ backgroundColor: '#f8f9fa' }}>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Order ID</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Amount</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Status</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Date</th>
                  <th style={{ padding: '1rem', textAlign: 'left', borderBottom: '1px solid #dee2e6' }}>Items</th>
                  <th style={{ padding: '1rem', textAlign: 'center', borderBottom: '1px solid #dee2e6' }}>Action</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((transaction, index) => {
                  const orderId = transaction.ID || transaction.id;
                  const status = transaction.Status || transaction.status;
                  const paymentUrl = transaction.PaymentURL || transaction.payment_url;

                  return (
                    <tr key={index} style={{ borderBottom: '1px solid #dee2e6' }}>
                      <td style={{ padding: '1rem' }}>{orderId}</td>
                      <td style={{ padding: '1rem' }}>Rp {transaction.Amount?.toLocaleString() || transaction.amount?.toLocaleString()}</td>
                      <td style={{ padding: '1rem' }}>
                        <span style={{ 
                          padding: '0.25rem 0.5rem',
                          borderRadius: '4px',
                          backgroundColor: getStatusColor(status),
                          color: 'white',
                          fontSize: '0.875rem'
                        }}>
                          {status || 'Unknown'}
                        </span>
                      </td>
                      <td style={{ padding: '1rem' }}>
                        {transaction.CreatedAt ? new Date(transaction.CreatedAt).toLocaleDateString() : 
                         transaction.created_at ? new Date(transaction.created_at).toLocaleDateString() : '-'}
                      </td>
                      <td style={{ padding: '1rem' }}>
                        {(transaction.Items || transaction.items)?.map((item, itemIndex) => (
                          <div key={itemIndex} style={{ fontSize: '0.825rem' }}>
                            {item.Name || item.name} (x{item.Quantity || item.quantity})
                          </div>
                        )) || '-'}
                      </td>
                      <td style={{ padding: '1rem', textAlign: 'center' }}>
                        <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
                          {status?.toLowerCase() === 'pending' && (
                            <>
                              <button 
                                onClick={() => handlePayNow(paymentUrl)}
                                style={{
                                  padding: '0.4rem 0.8rem',
                                  backgroundColor: '#007bff',
                                  color: 'white',
                                  border: 'none',
                                  borderRadius: '4px',
                                  cursor: 'pointer',
                                  fontSize: '0.8rem'
                                }}
                              >
                                Pay Now
                              </button>
                              <button 
                                onClick={() => handleCheckStatus(orderId)}
                                style={{
                                  padding: '0.4rem 0.8rem',
                                  backgroundColor: '#6c757d',
                                  color: 'white',
                                  border: 'none',
                                  borderRadius: '4px',
                                  cursor: 'pointer',
                                  fontSize: '0.8rem'
                                }}
                              >
                                Check
                              </button>
                            </>
                          )}
                          {status?.toLowerCase() === 'success' && <span style={{ color: '#28a745' }}>Completed</span>}
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default History;