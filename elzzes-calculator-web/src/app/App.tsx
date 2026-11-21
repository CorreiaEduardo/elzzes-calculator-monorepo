import { Calculator } from '../features/calculator/components/Calculator'
import { QueryProvider } from './QueryProvider'

export function App() {
  return (
    <QueryProvider>
      <main className="flex min-h-screen items-center justify-center bg-ui-page px-4 py-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] sm:p-6">
        <Calculator />
      </main>
    </QueryProvider>
  )
}
