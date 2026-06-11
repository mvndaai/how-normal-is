import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Supabase.initialize(
    url: const String.fromEnvironment('SUPABASE_URL'),
    anonKey: const String.fromEnvironment('SUPABASE_ANON_KEY'),
  );
  runApp(const HowNormalIsApp());
}

class HowNormalIsApp extends StatelessWidget {
  const HowNormalIsApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'How Normal Is',
      home: Supabase.instance.client.auth.currentSession == null
          ? const LoginPage()
          : const QuestionPage(),
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailController = TextEditingController();

  Future<void> _login() async {
    await Supabase.instance.client.auth.signInWithOtp(
      email: _emailController.text.trim(),
      emailRedirectTo: 'io.supabase.hownormalis://login-callback/',
    );
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Check your email for the login link.')),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('How Normal Is - Login')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(controller: _emailController, decoration: const InputDecoration(labelText: 'Email')),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _login, child: const Text('Login with magic link')),
          ],
        ),
      ),
    );
  }
}

class QuestionPage extends StatefulWidget {
  const QuestionPage({super.key});

  @override
  State<QuestionPage> createState() => _QuestionPageState();
}

class _QuestionPageState extends State<QuestionPage> {
  final _traitController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('How Normal Is')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Random unanswered question:'),
            const SizedBox(height: 8),
            const Text('How normal is it to prefer working at night?'),
            const SizedBox(height: 16),
            Wrap(
              spacing: 8,
              children: const [
                _AnswerButton(label: 'yes for me'),
                _AnswerButton(label: 'sometimes for me'),
                _AnswerButton(label: 'never for me'),
              ],
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _traitController,
              decoration: const InputDecoration(
                labelText: 'I think it is because <trait> (optional)',
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _AnswerButton extends StatelessWidget {
  final String label;
  const _AnswerButton({required this.label});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(onPressed: () {}, child: Text(label));
  }
}
