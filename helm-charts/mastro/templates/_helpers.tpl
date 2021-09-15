{{/*
Expand the name of the chart.
*/}}
{{- define "mastro.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "mastro.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "mastro.catalogueName" -}}
{{- printf "%s-%s-%s" .Chart.Name .Chart.Version "catalogue" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.fullCatalogueName" -}}
{{- printf "%s-%s" .Release.Name "catalogue" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.metricStoreName" -}}
{{- printf "%s-%s-%s" .Chart.Name .Chart.Version "metricstore" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.fullMetricStoreName" -}}
{{- printf "%s-%s" .Release.Name "metricstore" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.featureStoreName" -}}
{{- printf "%s-%s-%s" .Chart.Name .Chart.Version "featurestore" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.fullFeatureStoreName" -}}
{{- printf "%s-%s" .Release.Name "featurestore" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.uiName" -}}
{{- printf "%s-%s-%s" .Chart.Name .Chart.Version "ui" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "mastro.fullUIName" -}}
{{- printf "%s-%s" .Release.Name "ui" | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mastro.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mastro.labels" -}}
helm.sh/chart: {{ include "mastro.chart" . }}
{{ include "mastro.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mastro.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mastro.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}



{{/*
Create the name of the service account to use
*/}}
{{- define "mastro.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "mastro.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
