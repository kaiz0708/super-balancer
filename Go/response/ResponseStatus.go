package response

import (
	"Go/config"
	"html/template"
	"log"
	"net/http"
)

type MetricsPageData struct {
	Backends  map[string]*config.BackendMetrics
	Algorithm string
}

func HandleStatusHTML(w http.ResponseWriter, r *http.Request) {
	data := MetricsPageData{
		Backends:  config.MetricsMap,
		Algorithm: config.LoadBalancerDefault,
	}

	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Load Balancer Metrics</title>
    <style>
        :root {
            --primary-color: #2c3e50;
            --accent-color: #3498db;
            --healthy-color: #27ae60;
            --unhealthy-color: #e74c3c;
            --background-color: #ecf0f1;
            --card-background: #ffffff;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Roboto, sans-serif;
            background-color: var(--background-color);
            color: var(--primary-color);
            line-height: 1.6;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: var(--card-background);
            border-radius: 12px;
            padding: 30px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
        }

        h1 {
            text-align: center;
            color: var(--primary-color);
            font-size: 2.5rem;
            margin-bottom: 30px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .table-container {
            overflow-x: auto;
            border-radius: 8px;
            margin-top: 20px;
        }

        table {
            width: 100%;
            border-collapse: separate;
            border-spacing: 0;
            background: var(--card-background);
        }

        th, td {
            padding: 15px;
            text-align: left;
            border-bottom: 1px solid #e0e0e0;
            font-size: 0.95rem;
        }

        th {
            background-color: var(--accent-color);
            color: white;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            position: sticky;
            top: 0;
            z-index: 10;
        }

        tr {
            transition: background-color 0.3s ease;
        }

        tr:nth-child(even) {
            background-color: #f8f9fa;
        }

        tr:hover {
            background-color: #e8ecef;
        }

        .status {
            padding: 8px 12px;
            border-radius: 12px;
            font-weight: 500;
            text-align: center;
            display: inline-block;
        }

        .healthy {
            background-color: var(--healthy-color);
            color: white;
        }

        .unhealthy {
            background-color: var(--unhealthy-color);
            color: white;
        }

        .metric-number {
            font-family: 'Courier New', monospace;
            font-weight: 500;
        }

        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: var(--accent-color);
            color: white;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 500;
            transition: background-color 0.3s ease, transform 0.2s ease;
            margin: 20px 0;
            text-align: center;
        }

        .button:hover {
            background-color: #2980b9;
            transform: translateY(-2px);
        }

        .form-container {
            display: flex;
            gap: 10px;
            justify-content: center;
            margin: 20px 0;
        }

        .select-algorithm {
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 1rem;
            background: white;
            cursor: pointer;
        }

        .submit-button {
            padding: 12px 24px;
            background-color: var(--accent-color);
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 1rem;
            transition: background-color 0.3s ease;
        }

        .submit-button:hover {
            background-color: #2980b9;
        }

        .feedback-message {
            margin-top: 10px;
            text-align: center;
            font-size: 0.9rem;
            padding: 10px;
            border-radius: 6px;
            display: none;
        }

        .success {
            background-color: var(--healthy-color);
            color: white;
        }

        .error {
            background-color: var(--unhealthy-color);
            color: white;
        }

        @media (max-width: 768px) {
            .container {
                padding: 15px;
            }

            h1 {
                font-size: 1.8rem;
            }

            th, td {
                font-size: 0.85rem;
                padding: 10px;
            }

            .form-container {
                flex-direction: column;
                align-items: center;
            }

            .select-algorithm,
            .submit-button {
                width: 100%;
                max-width: 300px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Load Balancer Metrics</h1>
        <div class="table-container">
            <table>
                <tr>
                    <th>Backend</th>
                    <th>Request Count</th>
                    <th>Success Count</th>
                    <th>Failure Count</th>
                    <th>Total Latency</th>
                    <th>Avg Latency</th>
                    <th>Consecutive Fails</th>
                    <th>Consecutive Success</th>
                    <th>Timeout Break</th>
                    <th>Health Status</th>
                    <th>Last Status</th>
                    <th>Active Connections</th>
                    <th>Weight</th>
                    <th>Current Weight</th>
                </tr>
                {{range $key, $backend := .Backends}}
                <tr>
                    <td>{{$key}}</td>
                    <td class="metric-number">{{$backend.Metrics.RequestCount}}</td>
                    <td class="metric-number">{{$backend.Metrics.SuccessCount}}</td>
                    <td class="metric-number">{{$backend.Metrics.FailureCount}}</td>
                    <td class="metric-number">{{$backend.Metrics.TotalLatency}}</td>
                    <td class="metric-number">{{$backend.Metrics.AvgLatency}}</td>
                    <td class="metric-number">{{$backend.Metrics.ConsecutiveFails}}</td>
                    <td class="metric-number">{{$backend.Metrics.ConsecutiveSuccess}}</td>
                    <td class="metric-number">{{$backend.Metrics.TimeoutBreak}}</td>
                    <td>
                        <span class="status {{if $backend.Metrics.IsHealthy}}healthy{{else}}unhealthy{{end}}">
                            {{if $backend.Metrics.IsHealthy}}Healthy{{else}}Unhealthy{{end}}
                        </span>
                    </td>
                    <td>{{$backend.Metrics.LastStatus}}</td>
                    <td class="metric-number">{{$backend.Metrics.ActiveConnections}}</td>
                    <td class="metric-number">{{$backend.Metrics.Weight}}</td>
                    <td class="metric-number">{{$backend.Metrics.CurrentWeight}}</td>
                </tr>
                {{end}}
            </table>
        </div>
        <div class="form-container">
            <form id="change-algorithm-form">
                <select class="select-algorithm" name="name">
                    <option value="ROUND_ROBIN" {{if eq .Algorithm "ROUND_ROBIN"}}selected{{end}}>Round Robin</option>
					<option value="LEAST_CONNECTION" {{if eq .Algorithm "LEAST_CONNECTION"}}selected{{end}}>Least Connections</option>
					<option value="WEIGHTED_LEAST_CONNECTION" {{if eq .Algorithm "WEIGHTED_LEAST_CONNECTION"}}selected{{end}}>Weighted Least Connection</option>
					<option value="WEIGHTED_ROUND_ROBIN" {{if eq .Algorithm "WEIGHTED_ROUND_ROBIN"}}selected{{end}}>Weighted Round Robin</option>
					<option value="RANDOM" {{if eq .Algorithm "RANDOM"}}selected{{end}}>Random</option>
                </select>
                <button type="submit" class="submit-button">Apply Algorithm</button>
            </form>
        </div>
        <div id="feedback-message" class="feedback-message"></div>
    </div>
    <script>
        document.getElementById('change-algorithm-form').addEventListener('submit', async function(event) {
            event.preventDefault();
            
            const selectElement = this.querySelector('.select-algorithm');
            const algorithm = selectElement.value;
            const feedbackMessage = document.getElementById('feedback-message');

            try {
                const response = await fetch('/change-load-balancer?name=' + encodeURIComponent(algorithm), {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const result = await response.json();
                    feedbackMessage.textContent = 'Successfully changed to ' + result + ' algorithm';
                    feedbackMessage.classList.remove('error');
                    feedbackMessage.classList.add('success');
                } else {
                    throw new Error('Failed to change algorithm');
                }
            } catch (error) {
                feedbackMessage.textContent = 'Error: ' + error.message;
                feedbackMessage.classList.remove('success');
                feedbackMessage.classList.add('error');
            }

            feedbackMessage.style.display = 'block';
            setTimeout(() => {
                feedbackMessage.style.display = 'none';
            }, 3000);
        });
    </script>
</body>
</html>`

	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		log.Println("Error rendering HTML:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
