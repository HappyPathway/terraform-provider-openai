import pandas as pd
import matplotlib.pyplot as plt

def analyze_data(data_path):
    # Read the CSV file
    df = pd.read_csv(data_path)
    
    # Calculate basic statistics
    stats = {
        'total_revenue': df['revenue'].sum(),
        'total_expenses': df['expenses'].sum(),
        'total_profit': df['profit'].sum(),
        'avg_daily_profit': df['profit'].mean()
    }
    
    # Create a visualization
    plt.figure(figsize=(10, 6))
    plt.plot(df['date'], df['profit'], marker='o', label='Profit')
    plt.plot(df['date'], df['revenue'], marker='s', label='Revenue')
    plt.plot(df['date'], df['expenses'], marker='^', label='Expenses')
    plt.title('Daily Financial Performance')
    plt.xlabel('Date')
    plt.ylabel('Amount ($)')
    plt.legend()
    plt.xticks(rotation=45)
    plt.tight_layout()
    
    return stats